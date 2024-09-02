package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitlab.ozon.dev/chppppr/homework/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

func main() {
	storage := storage.NewStorage("storage.json")
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
