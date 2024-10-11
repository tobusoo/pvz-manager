package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
	manager_service "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
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

	ordrs := make([]uint64, 0, len(orders))
	for _, v := range orders {
		ordrs = append(ordrs, uint64(v))
	}

	req := &manager_service.GiveOrdersRequestV1{
		Orders: ordrs,
	}

	task := &workers.TaskRequest{
		Request: "give -o=...",
		Func: func() error {
			_, err := mng_service.GiveOrdersV1(ctx, req)
			return err
		},
	}

	fmt.Printf("\n\n")
	wk.AddTask(task)
}
