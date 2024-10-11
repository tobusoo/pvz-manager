package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	"gitlab.ozon.dev/chppppr/homework/internal/workers"
	manager_service "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	exp_date, err := utils.StringToTime(expirationDate)
	if err != nil {
		wk.Results <- &workers.TaskResponse{
			Response: request_str,
			Err:      err,
		}
		return
	}

	order := &manager_service.OrderV1{
		ExpirationDate: timestamppb.New(exp_date),
		PackageType:    containerType,
		UseTape:        useTape,
		Cost:           cost,
		Weight:         weight,
	}

	req := &manager_service.AddOrderRequestV1{
		OrderId: orderID,
		UserId:  userID,
		Order:   order,
	}

	task := &workers.TaskRequest{
		Request: request_str,
		Func: func() error {
			_, err := mng_service.AddOrderV1(ctx, req)
			return err
		},
	}

	fmt.Printf("\n\n")
	wk.AddTask(task)
}

func acceptRefundCmdRun(cmd *cobra.Command, args []string) {
	defer resetRefundFlags(cmd)

	req := &manager_service.RefundRequestV1{
		UserId:  userID,
		OrderId: orderID,
	}

	task := &workers.TaskRequest{
		Request: fmt.Sprintf("accept refund -u=%d -o=%d", userID, orderID),
		Func: func() error {
			_, err := mng_service.RefundV1(ctx, req)
			return err
		},
	}

	fmt.Printf("\n\n")
	wk.AddTask(task)
}
