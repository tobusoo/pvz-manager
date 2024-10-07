package app

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
)

func readInput(exit *bool, inputCh chan string) {
	reader := bufio.NewReader(os.Stdin)

	for !(*exit) {
		cmd.InOutLock()
		fmt.Print(">>> ")
		input, _ := reader.ReadString('\n')
		cmd.InOutUnlock()

		input = strings.TrimSpace(input)
		inputCh <- input

		time.Sleep(25 * time.Millisecond)
	}

	close(inputCh)
}

func processInput(exit *bool, input string) {
	if input == "exit" {
		*exit = true
		return
	}

	args := strings.Fields(input)
	if len(args) > 0 {
		cmd.SetArgs(args)
		cmd.Execute()
	}
}

func RunInteractive(ctx context.Context) {
	wk := workers.NewWorkers(1)
	cmd.SetWorker(wk)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go ShowResult(wg, wk)
	defer cmd.CloseAndWaitWorkers()

	exit := false
	inputCh := make(chan string)
	go readInput(&exit, inputCh)

	for !exit {
		select {
		case <-ctx.Done():
			exit = true
			fmt.Println()
			log.Println("Receive os signal")
			cmd.CloseAndWaitWorkers()
			log.Println("All workers finished their jobs")
		case input := <-inputCh:
			processInput(&exit, input)
		}
	}
}
