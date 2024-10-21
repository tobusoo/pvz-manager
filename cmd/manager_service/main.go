package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gitlab.ozon.dev/chppppr/homework/internal/app/manager_service"
	kafka_client "gitlab.ozon.dev/chppppr/homework/internal/clients/kafka"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka/producer"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
	"gitlab.ozon.dev/chppppr/homework/internal/usecase"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func init() {
	_ = godotenv.Load()
}

func newManagerService(ctx context.Context, pool *pgxpool.Pool, pr sarama.SyncProducer, cfg *Config) (*manager_service.ManagerService, error) {
	txManager := postgres.NewTxManager(pool)
	pgPepo := postgres.NewRepoPG(txManager)
	st := postgres.NewStorageDB(ctx, txManager, pgPepo)
	au := usecase.NewAcceptUsecase(st)
	gu := usecase.NewGiveUsecase(st)
	ru := usecase.NewReturnUsecase(st)
	vu := usecase.NewViewUsecase(st)

	pr_client := kafka_client.NewProducerClient(pr, cfg.Kafka.Topic)
	return manager_service.NewManagerService(au, gu, ru, vu, pr_client), nil
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctxWichCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	pool, err := pgxpool.New(ctxWichCancel, os.Getenv("POSTGRESQL_DSN"))
	if err != nil {
		log.Fatal("pgxpool.New:", err)
	}
	defer pool.Close()

	pr, err := producer.NewSyncProducer(cfg.Kafka.Config)
	if err != nil {
		log.Fatal("producer.NewSyncProducer:", err)
	}
	defer pr.Close()

	mng_service, err := newManagerService(ctxWichCancel, pool, pr, cfg)
	if err != nil {
		log.Fatal("newManagerService:", err)
	}

	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	desc.RegisterManagerServiceServer(grpcServer, mng_service)

	mux := runtime.NewServeMux()
	err = desc.RegisterManagerServiceHandlerFromEndpoint(ctxWichCancel, mux, cfg.GRPC.Address, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatalf("failed to register manager service handler: %v", err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/v1/", mux)
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/manager-service/v1/manager-service.swagger.json")
	})
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://"+cfg.Swagger.Address+"/swagger.json")))

	httpServer := http.Server{Addr: cfg.HTPP.Address, Handler: r}
	go httpServer.ListenAndServe()

	<-ctxWichCancel.Done()
	fmt.Println()
	log.Println("Receive os signal")
	grpcServer.GracefulStop()
	httpServer.Shutdown(context.Background())
	log.Println("all done")
}
