package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	viewCmd.AddCommand(viewRefundCmd)
}

var (
	viewCmd = &cobra.Command{
		Use:   "view",
		Short: "View orders or refunds",
		Long:  "View orders or refunds",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	viewRefundCmd = &cobra.Command{
		Use:   "refund",
		Short: "View refunds",
		Long:  "View refunds",
		Run:   viewRefundCmdRun,
	}
)

func viewRefundCmdRun(cmd *cobra.Command, args []string) {
	refunds, err := st.GetRefunds()
	if err != nil {
		fmt.Printf("error while view refund: %s\n", err)
		return
	}

	if len(refunds) == 0 {
		fmt.Println("there are no refunds")
		return
	}

	for _, orderID := range refunds {
		fmt.Printf("orderID %d\n", orderID)
	}
}
