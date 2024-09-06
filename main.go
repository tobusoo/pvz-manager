package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gitlab.ozon.dev/chppppr/homework/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

func main() {
	storage, err := storage.NewStorage("storage.json")
	if err != nil {
		fmt.Println(err)
	}
	defer storage.Save()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		storage.Save()
		fmt.Println()
		os.Exit(0)
	}()

	cmd.SetStorage(storage)
	if len(os.Args[1:]) > 0 {
		if err := cmd.Execute(); err != nil {
			fmt.Println(err)
		}
		return
	}

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
