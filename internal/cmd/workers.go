package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
)

const MaxWorkers = 100
const MinWorkers = 1

func init() {
	workersCmd.AddCommand(workersSetCmd)
	workersCmd.AddCommand(workersViewCmd)

	resetWorkersFlags(workersSetCmd)
	workersSetCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetWorkersFlags(cmd)
	})

	workersViewCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	})
}

var (
	workersCmd = &cobra.Command{
		Use:   "workers",
		Short: "Set or view num of workers",
		Long:  "Set or view num of workers",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	workersSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set num of workers",
		Long:  "Set num of workers",
		Run:   workersSetCmdRun,
	}

	workersViewCmd = &cobra.Command{
		Use:   "num",
		Short: "View num of workers",
		Long:  "View num of workers",
		Run:   workersViewCmdRun,
	}
)

func resetWorkersFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().UintVarP(&numWorkers, "num", "n", 1, "Num of workers (required)")
	cmd.MarkPersistentFlagRequired("num")
}

func workersSetCmdRun(cmd *cobra.Command, args []string) {
	defer resetWorkersFlags(cmd)

	if numWorkers > MaxWorkers {
		wk.Results <- &workers.TaskResponse{
			Err:      fmt.Errorf("MaxWorkes = %d", MaxWorkers),
			Response: fmt.Sprintf("workers -n %d", numWorkers)}
		return
	}

	if numWorkers < MinWorkers {
		wk.Results <- &workers.TaskResponse{
			Err:      fmt.Errorf("MinWorkers = %d", MinWorkers),
			Response: fmt.Sprintf("workers -n %d", numWorkers)}
		return
	}

	CloseAndWaitWorkers()
	SetWorkers(numWorkers)
}

func workersViewCmdRun(cmd *cobra.Command, args []string) {
	InOutLock()
	fmt.Println("Workers num: ", wk.GetSize())
	InOutUnlock()
}
