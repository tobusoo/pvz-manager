package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/storage_json"
	"gitlab.ozon.dev/chppppr/homework/internal/usecase"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
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

	wk := workers.NewWorkers(2)
	cmd.SetWorkers(wk)
	acceptUsecase := usecase.NewAcceptUsecase(st)

	result_handler := func(wk *workers.Workers) {
		for res := range wk.Results {
			fmt.Println(res)
		}
	}

	go result_handler(wk)
	for i := 0; i < 50; i++ {
		req := &dto.AddOrderRequest{
			ExpirationDate: "10-10-2024",
			ContainerType:  "",
			UseTape:        false,
			Cost:           100,
			Weight:         100,
			OrderID:        uint64(1 + i),
			UserID:         1,
		}
		if i == 25 {
			fmt.Println("Set new workers")
			wk.Close()
			wk.Wait()
			wk = workers.NewWorkers(5)
			go result_handler(wk)
		}
		task := workers.TaskRequest{Func: func() error {
			return acceptUsecase.AcceptOrder(req)
		}, Request: strconv.Itoa(i)}

		wk.AddTask(task)
		fmt.Println("add", i)
	}
	wk.Close()
	wk.Wait()

	// if len(os.Args[1:]) > 0 {
	// 	RunOnce()
	// } else {
	// 	RunInteractive()
	// }
}
