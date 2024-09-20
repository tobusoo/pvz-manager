package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

func init() {
	_ = godotenv.Load()
}

func RunWithExit(st *storage.Storage) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}

	st.Save()
	os.Exit(0)
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
	ordersHistoryRep := storage.NewOrdersHistory()
	refundsRep := storage.NewRefunds()
	usersRep := storage.NewUsers()
	storagePath := "storage.json"
	if envStoragePath, ok := os.LookupEnv("STORAGE_PATH"); ok {
		storagePath = envStoragePath
	}

	storage, err := storage.NewStorage(ordersHistoryRep, refundsRep, usersRep, storagePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer storage.Save()

	cmd.SetStorage(storage)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		storage.Save()
		fmt.Println()
		os.Exit(0)
	}()

	if len(os.Args[1:]) > 0 {
		RunWithExit(storage)
	}

	RunInteractive()
}
