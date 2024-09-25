package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
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
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	txManager := postgres.NewTxManager(pool)
	pgPepo := postgres.NewRepoPG(txManager)
	storage := postgres.NewStorageDB(ctx, txManager, pgPepo)

	cmd.SetStorage(storage)

	if len(os.Args[1:]) > 0 {
		RunOnce()
	} else {
		RunInteractive()
	}
}
