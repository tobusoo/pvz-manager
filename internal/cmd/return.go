package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
)

func init() {
	resetReturnFlags(returnCmd)
	returnCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetReturnFlags(cmd)
	})
}

var returnCmd = &cobra.Command{
	Use:   "return",
	Short: "Return order to courier",
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

	req := &dto.ReturnRequest{
		OrderID: orderID,
	}

	if err := returnUsecase.Return(req); err != nil {
		fmt.Println(err)
	}
}
