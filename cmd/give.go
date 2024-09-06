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

	knowUserID := false
	userID := uint64(0)

	slices.Sort(orders)
	orders = slices.Compact(orders)

	errors := make([]error, 0)
	isGoodResponse := true

	for _, order := range orders {
		status, err := st.GetOrderStatus(uint64(order))
		if err != nil {
			errors = append(errors, fmt.Errorf("can't give: %s", err))
			isGoodResponse = false
			continue
		}

		if !knowUserID {
			userID = status.UserID
			knowUserID = true
		}

		if status.UserID != userID {
			errors = append(errors, fmt.Errorf("can't give order %d: different userID", order))
			isGoodResponse = false
			continue
		}

		if status.Status != storage.StatusAccepted {
			errors = append(errors, fmt.Errorf("can't give order %d: status = %s", order, status.Status))
			isGoodResponse = false
			continue
		}

		expDate, err := st.GetExpirationDate(status.UserID, uint64(order))
		if err != nil {
			errors = append(errors, fmt.Errorf("can't give: %s", err))
			isGoodResponse = false
			continue
		}

		if utils.CurrentDate().After(expDate) {
			errors = append(errors, fmt.Errorf("can't give order %d: expiration date has already passed", order))
			isGoodResponse = false
			continue
		}

		if err = st.CanRemoveOrder(uint64(order)); err != nil {
			errors = append(errors, err)
			isGoodResponse = false
		}
	}

	if isGoodResponse {
		for _, order := range orders {
			st.RemoveOrder(uint64(order), storage.StatusGiveClient)
		}
	} else {
		fmt.Println("request was not done because:")
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}
