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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gitlab.ozon.dev/chppppr/homework/internal/app/manager_service"
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

func newManagerService(ctx context.Context, pool *pgxpool.Pool) *manager_service.ManagerService {
	txManager := postgres.NewTxManager(pool)
	pgPepo := postgres.NewRepoPG(txManager)
	st := postgres.NewStorageDB(ctx, txManager, pgPepo)
	au := usecase.NewAcceptUsecase(st)
	gu := usecase.NewGiveUsecase(st)
	ru := usecase.NewReturnUsecase(st)
	vu := usecase.NewViewUsecase(st)

	return manager_service.NewManagerService(au, gu, ru, vu)
}

func main() {
	ctx := context.Background()
	ctxWichCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	pool, err := pgxpool.New(ctxWichCancel, os.Getenv("POSTGRESQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	mng_service := newManagerService(ctxWichCancel, pool)

	lis, err := net.Listen("tcp", os.Getenv("GRPC_HOST"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	desc.RegisterManagerServiceServer(grpcServer, mng_service)

	mux := runtime.NewServeMux()
	err = desc.RegisterManagerServiceHandlerFromEndpoint(ctxWichCancel, mux, os.Getenv("GRPC_HOST"), []grpc.DialOption{
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
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://"+os.Getenv("HTTP_HOST")+"/swagger.json")))

	httpServer := http.Server{Addr: os.Getenv("HTTP_HOST"), Handler: r}
	go httpServer.ListenAndServe()

	<-ctxWichCancel.Done()
	fmt.Println()
	log.Println("Receive os signal")
	grpcServer.GracefulStop()
	httpServer.Shutdown(context.Background())
	log.Println("all done")
}
