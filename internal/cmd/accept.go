package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
)

func init() {
	acceptCmd.AddCommand(acceptOrderCmd)
	acceptCmd.AddCommand(acceptRefundCmd)

	resetOrderFlags(acceptOrderCmd)
	acceptOrderCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetOrderFlags(cmd)
	})

	resetRefundFlags(acceptRefundCmd)
	acceptRefundCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetRefundFlags(cmd)
	})
}

var (
	acceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "Accept orders or refund",
		Long:  "Accept orders or refund",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	acceptOrderCmd = &cobra.Command{
		Use:   "order",
		Short: "Accept order",
		Long:  "Accept order",
		Run:   acceptOrderCmdRun,
	}

	acceptRefundCmd = &cobra.Command{
		Use:   "refund",
		Short: "Accept refund",
		Long:  "Accept refund",
		Run:   acceptRefundCmdRun,
	}
)

func resetRefundFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().Uint64VarP(&userID, "userID", "u", 0, "userID (required)")
	cmd.PersistentFlags().Uint64VarP(&orderID, "orderID", "o", 0, "orderID (required)")
	cmd.MarkPersistentFlagRequired("userID")
	cmd.MarkPersistentFlagRequired("orderID")
}

func resetOrderFlags(cmd *cobra.Command) {
	resetRefundFlags(cmd)

	cmd.PersistentFlags().Uint64VarP(&cost, "cost", "c", 0, "cost in rubles (required)")
	cmd.PersistentFlags().Uint64VarP(&weight, "weight", "w", 0, "weight in grams (required)")
	cmd.PersistentFlags().StringVarP(&expirationDate, "time", "t", "", "Expiration Date (required)")
	cmd.MarkPersistentFlagRequired("cost")
	cmd.MarkPersistentFlagRequired("weight")
	cmd.MarkPersistentFlagRequired("time")

	cmd.PersistentFlags().StringVarP(&containerType, "containerType", "p", "", "containerType (tape, package, box)")
	cmd.PersistentFlags().BoolVarP(&useTape, "useTape", "s", false, "use additional tape (containerType must be defined)")
}

func acceptOrderCmdRun(cmd *cobra.Command, args []string) {
	defer resetOrderFlags(cmd)

	request_str := fmt.Sprintf("accept order -u=%d -o=%d ...", userID, orderID)
	req := &dto.AddOrderRequest{
		UserID:         userID,
		OrderID:        orderID,
		ExpirationDate: expirationDate,
		ContainerType:  containerType,
		UseTape:        useTape,
		Cost:           cost,
		Weight:         weight,
	}

	task := &workers.TaskRequest{
		Request: request_str,
		Func: func() error {
			return mng_client.AddOrder(ctx, req)
		},
	}

	fmt.Printf("\n\n")
	wk.AddTask(task)
}

func acceptRefundCmdRun(cmd *cobra.Command, args []string) {
	defer resetRefundFlags(cmd)

	req := &dto.RefundRequest{
		UserID:  userID,
		OrderID: orderID,
	}

	task := &workers.TaskRequest{
		Request: fmt.Sprintf("accept refund -u=%d -o=%d", userID, orderID),
		Func: func() error {
			return mng_client.Refund(ctx, req)
		},
	}

	fmt.Printf("\n\n")
	wk.AddTask(task)
}
