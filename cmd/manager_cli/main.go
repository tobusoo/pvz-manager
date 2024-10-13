package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chppppr/homework/internal/app/manager_cli"
	"gitlab.ozon.dev/chppppr/homework/internal/clients/manager"
	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	manager_service "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	ctx := context.Background()
	ctxWichCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	conn, err := grpc.NewClient(os.Getenv("GRPC_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	mng_service := manager_service.NewManagerServiceClient(conn)
	mng_client := manager.NewManagerServiceClient(mng_service)
	cmd.SetManagerServiceClient(mng_client)
	cmd.SetContext(ctxWichCancel)

	if len(os.Args[1:]) > 0 {
		manager_cli.RunOnce()
	} else {
		manager_cli.RunInteractive(ctxWichCancel)
	}
}
