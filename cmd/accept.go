package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

func init() {
	acceptOrderCmd.DisableSuggestions = false
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
	cmd.PersistentFlags().StringVarP(&expirationDate, "time", "t", "", "Expiration Date (required)")
	cmd.MarkPersistentFlagRequired("time")
}

func acceptOrderCmdRun(cmd *cobra.Command, args []string) {
	defer resetOrderFlags(cmd)

	expDate, err := time.Parse("02-01-2006", expirationDate)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentTime := time.Now().Truncate(24 * time.Hour)
	if currentTime.After(expDate) {
		fmt.Println("expiration date has already passed")
		return
	}

	if err := st.AddOrder(userID, orderID, expDate.Format("02-01-2006")); err != nil {
		fmt.Println(err)
	}
}

func acceptRefundCmdRun(cmd *cobra.Command, args []string) {
	defer resetRefundFlags(cmd)

	order, err := st.GetOrderStatus(orderID)
	if err != nil {
		fmt.Println(err, st.OrdersHistory)
		return
	}

	if order.Status != storage.StatusGiveClient {
		fmt.Printf("can not refund order %d: status = %s\n", orderID, order.Status)
		return
	}

	if userID != order.UserID {
		fmt.Printf("can not refund order %d: wrong userID\n", orderID)
		return
	}

	issuedDate, err := time.Parse("02-01-2006", order.Date)
	if err != nil {
		fmt.Println(err)
		return
	}

	issuedDate = issuedDate.Add(2 * 24 * time.Hour)
	currentDate := time.Now().Truncate(24 * time.Hour)

	if currentDate.After(issuedDate) {
		fmt.Println("2 days have passed since the order was issued to the client")
		return
	}

	if err = st.AddRefund(orderID); err != nil {
		fmt.Println(err)
	}
}
