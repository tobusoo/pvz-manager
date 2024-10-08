package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
)

func init() {
	giveCmd.DisableSuggestions = false

	resetGiveCmd(giveCmd)
	giveCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetGiveCmd(cmd)
	})
}

var (
	giveCmd = &cobra.Command{
		Use:   "give",
		Short: "Give orders to client",
		Long:  "Give orders to client",
		Run:   giveCmdRun,
	}
)

func resetGiveCmd(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().UintSliceVarP(&orders, "orders", "o", []uint{}, "List of orderID")
	cmd.MarkPersistentFlagRequired("orders")
}

func giveCmdRun(cmd *cobra.Command, args []string) {
	defer resetGiveCmd(cmd)

	req := &dto.GiveOrdersRequest{
		Orders: orders,
	}

	task := &workers.TaskRequest{
		Request: "give -o=...",
		Func: func() error {
			errs := giveUsecase.Give(req)
			return errors.Join(errs...)
		},
	}

	fmt.Printf("\n\n")
	wk.AddTask(task)
}
