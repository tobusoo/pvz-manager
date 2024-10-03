package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chppppr/homework/internal/app"
	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/storage_json"
)

func init() {
	_ = godotenv.Load()
}

func newStorageJSON() *storage_json.Storage {
	ordersHistoryRep := storage_json.NewOrdersHistory()
	refundsRep := storage_json.NewRefunds()
	usersRep := storage_json.NewUsers()
	storagePath := "storage.json"
	if envStoragePath, ok := os.LookupEnv("STORAGE_PATH"); ok {
		storagePath = envStoragePath
	}

	st, err := storage_json.NewStorage(ordersHistoryRep, refundsRep, usersRep, storagePath)
	if err != nil {
		log.Fatal(err)
	}

	return st
}

func main() {
	var st storage.Storage

	ctx := context.Background()
	ctxWichCancel, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if yes, ok := os.LookupEnv("USE_POSTGRESQL"); ok && yes == "yes" {
		pool, err := pgxpool.New(ctxWichCancel, os.Getenv("POSTGRESQL_DSN"))
		if err != nil {
			log.Fatal(err)
		}
		defer pool.Close()

		txManager := postgres.NewTxManager(pool)
		pgPepo := postgres.NewRepoPG(txManager)
		st = postgres.NewStorageDB(ctxWichCancel, txManager, pgPepo)
	} else {
		storage := newStorageJSON()
		defer storage.Save()

		st = storage
	}
	cmd.SetStorage(st)

	if len(os.Args[1:]) > 0 {
		app.RunOnce()
	} else {
		app.RunInteractive(ctxWichCancel)
	}
}
