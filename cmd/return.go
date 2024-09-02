package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

func init() {
	returnCmd.DisableSuggestions = false

	resetReturnFlags(returnCmd)
	returnCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetReturnFlags(cmd)
	})
}

var returnCmd = &cobra.Command{
	Use:   "return",
	Short: "Return order",
	Long:  "Return order to courier",
	Run:   returnCmdRun,
}

func resetReturnFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().Uint64VarP(&orderID, "orderID", "o", 0, "orderID (required)")
	cmd.MarkPersistentFlagRequired("orderID")
}

func returnCmdRun(cmd *cobra.Command, args []string) {
	defer resetReturnFlags(cmd)

	order, err := st.GetOrderStatus(orderID)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch order.Status {
	case storage.StatusReturned:
		st.RemoveReturned(orderID)
		st.SetOrderStatus(orderID, storage.StatusGiveCourier)

	case storage.StatusAccepted:
		expDate, err := st.GetExpirationDate(order.UserID, orderID)
		if err != nil {
			fmt.Println(err)
			return
		}

		if expDate.Before(utils.CurrentDate()) {
			fmt.Printf("can't return order %d: expiration date hasn't expired yet\n", orderID)
			return
		}

		if err := st.RemoveOrder(orderID, storage.StatusGiveCourier); err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Printf("can't return order %d: status = %s\n", orderID, order.Status)
	}
}
