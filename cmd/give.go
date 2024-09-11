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

func giveCheckErr(userID, orderID uint64, status storage.OrderStatus) error {
	if status.UserID != userID {
		return fmt.Errorf("can't give order %d: different userID", orderID)
	}

	if status.Status != storage.StatusAccepted {
		return fmt.Errorf("can't give order %d: status = %s", orderID, status.Status)
	}

	expDate, err := st.GetExpirationDate(status.UserID, uint64(orderID))
	if err != nil {
		return fmt.Errorf("can't give: %s", err)
	}

	if utils.CurrentDate().After(expDate) {
		return fmt.Errorf("can't give order %d: expiration date has already passed", orderID)
	}

	return st.CanRemoveOrder(orderID)
}

func giveCmdProcess(isGoodResponse bool, errors []error) {
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

func giveCheckOrder(orderID, userID uint64, knowUserID bool) (uint64, bool, error) {
	status, err := st.GetOrderStatus(orderID)
	if err != nil {
		return userID, knowUserID, fmt.Errorf("can't give: %s", err)
	}

	if !knowUserID {
		userID = status.UserID
		knowUserID = true
	}

	return userID, knowUserID, giveCheckErr(userID, orderID, status)
}

func giveCmdRun(cmd *cobra.Command, args []string) {
	defer resetGiveCmd(cmd)

	var err error
	userID := uint64(0)
	knowUserID := false
	isGoodResponse := true
	errors := make([]error, 0)

	slices.Sort(orders)
	orders = slices.Compact(orders)

	for _, orderID := range orders {
		userID, knowUserID, err = giveCheckOrder(uint64(orderID), userID, knowUserID)
		if err != nil {
			errors = append(errors, err)
			isGoodResponse = false
		}
	}

	giveCmdProcess(isGoodResponse, errors)
}
