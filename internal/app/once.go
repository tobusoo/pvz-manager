package app

import (
	"fmt"
	"sync"

	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
)

func RunOnce() {
	cmd.SetWorkers(1)
	defer cmd.CloseAndWaitWorkers()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}

	wg := &sync.WaitGroup{}
	getWorkersAndShowResult(wg)
	wg.Wait()
}
