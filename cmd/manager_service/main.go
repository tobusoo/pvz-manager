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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chppppr/homework/internal/app/manager_service"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	var st storage.Storage

	ctx := context.Background()
	ctxWichCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	pool, err := pgxpool.New(ctxWichCancel, os.Getenv("POSTGRESQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	txManager := postgres.NewTxManager(pool)
	pgPepo := postgres.NewRepoPG(txManager)
	st = postgres.NewStorageDB(ctxWichCancel, txManager, pgPepo)
	mng_service := manager_service.NewManagerService(st)

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
		if err := http.ListenAndServe(os.Getenv("HTTP_HOST"), mux); err != nil {
			log.Fatalf("failed to listen and serve manager service handler: %v", err)
		}
	}()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctxWichCancel.Done()
	fmt.Println()
	log.Println("Receive os signal")
	grpcServer.GracefulStop()
	log.Println("all done")
	// TODO: http gateway graceful stop
}
