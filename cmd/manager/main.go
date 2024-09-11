package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

func RunWithExit() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
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
	storage, err := storage.NewStorage("storage.json")
	if err != nil {
		fmt.Println(err)
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
		RunWithExit()
	}

	RunInteractive()
}
