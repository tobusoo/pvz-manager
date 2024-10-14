package cmd

import (
	"context"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/clients"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
)

func init() {
	rootCmd.AddCommand(acceptCmd)
	rootCmd.AddCommand(giveCmd)
	rootCmd.AddCommand(returnCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(workersCmd)

	rootCmd.DisableSuggestions = false
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if len(os.Args) == 1 {
		rootCmd.Use = ""
	}
}

var (
	inoutMtx       sync.Mutex
	numWorkers     uint
	prevNumWorkers uint
	wk             *workers.Workers

	mng_client clients.ManagerService
	ctx        context.Context

	cost           uint64
	weight         uint64
	pageID         uint64
	userID         uint64
	orderID        uint64
	useTape        bool
	ordersLimit    uint64
	ordersPerPage  uint64
	containerType  string
	expirationDate string
	orders         []uint

	rootCmd = &cobra.Command{
		Use:  "manager",
		Long: `Utility for managing an order pick-up point`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func SetManagerServiceClient(mng clients.ManagerService) {
	mng_client = mng
}

func SetContext(context context.Context) {
	ctx = context
}

func SetWorker(workers *workers.Workers) {
	wk = workers
	prevNumWorkers = wk.GetSize()
}

func SetWorkersNum(num uint) {
	if prevNumWorkers < numWorkers {
		wk.AddWorkers(numWorkers - prevNumWorkers)
	} else if prevNumWorkers > numWorkers {
		wk.CloseNworkers(prevNumWorkers - numWorkers)
	}
	prevNumWorkers = numWorkers
}

func GetWorker() *workers.Workers {
	return wk
}

func InOutLock() {
	inoutMtx.Lock()
}

func InOutUnlock() {
	inoutMtx.Unlock()
}

func CloseAndWaitWorkers() {
	wk.CloseJobs()
	wk.Wait()
}

func SetArgs(args []string) {
	rootCmd.SetArgs(args)
}

func Execute() error {
	return rootCmd.Execute()
}
