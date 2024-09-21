package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
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

	if errors := giveUsecase.Give(req); errors != nil {
		fmt.Println("request was not done because:")
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}
