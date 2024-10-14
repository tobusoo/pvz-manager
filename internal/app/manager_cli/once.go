package manager_cli

import (
	"fmt"
	"sync"

	"gitlab.ozon.dev/chppppr/homework/internal/cmd"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
)

func RunOnce() {
	wk := workers.NewWorkers(1)
	cmd.SetWorker(wk)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go ShowResult(wg, wk)
	defer cmd.CloseAndWaitWorkers()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
