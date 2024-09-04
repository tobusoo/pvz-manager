package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
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
	orders []uint

	giveCmd = &cobra.Command{
		Use:   "give",
		Short: "Give orders to client",
		Long:  "Give orders to client",
		Run:   fulfillCmdRun,
	}
)

func resetGiveCmd(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().UintSliceVarP(&orders, "orders", "o", []uint{}, "List of orderID")
	cmd.MarkPersistentFlagRequired("orders")
}

func fulfillCmdRun(cmd *cobra.Command, args []string) {
	defer resetGiveCmd(cmd)

	knowUserID := false
	userID := uint64(0)

	slices.Sort(orders)
	orders = slices.Compact(orders)

	for _, order := range orders {
		status, err := st.GetOrderStatus(uint64(order))
		if err != nil {
			fmt.Printf("can't give: %s\n", err)
			continue
		}

		if !knowUserID {
			userID = status.UserID
			knowUserID = true
		}

		if status.UserID != userID {
			fmt.Printf("can't give order %d: different userID", order)
			continue
		}

		if status.Status != storage.StatusAccepted {
			fmt.Printf("can't give order %d: status = %s\n", order, status.Status)
			continue
		}

		expDate, err := st.GetExpirationDate(status.UserID, uint64(order))
		if err != nil {
			fmt.Printf("can't give: %s\n", err)
			continue
		}

		if utils.CurrentDate().After(expDate) {
			fmt.Printf("can't give order %d: expiration date has already passed\n", order)
			continue
		}

		if err = st.RemoveOrder(uint64(order), storage.StatusGiveClient); err != nil {
			fmt.Println(err)
		}
	}
}
