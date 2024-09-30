package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/storage_json"
)

func init() {
	_ = godotenv.Load()
}

func RunOnce() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func RunInteractive() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">>> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			return
		}
		args := strings.Fields(input)

		if len(args) > 0 {
			cmd.SetArgs(args)
			cmd.Execute()
		}
	}
}

func main() {
	var st storage.Storage

	if yes, ok := os.LookupEnv("USE_POSTGRESQL"); ok && yes == "yes" {
		ctx := context.Background()

		pool, err := pgxpool.New(ctx, os.Getenv("POSTGRESQL_DSN"))
		if err != nil {
			log.Fatal(err)
		}
		defer pool.Close()

		txManager := postgres.NewTxManager(pool)
		pgPepo := postgres.NewRepoPG(txManager)
		st = postgres.NewStorageDB(ctx, txManager, pgPepo)
	} else {
		ordersHistoryRep := storage_json.NewOrdersHistory()
		refundsRep := storage_json.NewRefunds()
		usersRep := storage_json.NewUsers()
		storagePath := "storage.json"
		if envStoragePath, ok := os.LookupEnv("STORAGE_PATH"); ok {
			storagePath = envStoragePath
		}

		storage, err := storage_json.NewStorage(ordersHistoryRep, refundsRep, usersRep, storagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer storage.Save()

		st = storage
	}

	cmd.SetStorage(st)

	if len(os.Args[1:]) > 0 {
		RunOnce()
	} else {
		RunInteractive()
	}
}
