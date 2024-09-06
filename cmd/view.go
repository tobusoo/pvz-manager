package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	viewCmd.AddCommand(viewRefundCmd)
	viewCmd.AddCommand(viewOrderCmd)
	viewCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		cmd.ResetFlags()
	})

	resetViewOrderFlags(viewOrderCmd)
	viewOrderCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetViewOrderFlags(cmd)
	})

	resetViewRefundFlags(viewRefundCmd)
	viewRefundCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetViewRefundFlags(cmd)
	})
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

	pageID, ordersPerPage uint64
	viewRefundCmd         = &cobra.Command{
		Use:   "refund",
		Short: "View refunds",
		Long:  "View refunds",
		Run:   viewRefundCmdRun,
	}

	ordersLimit  uint64
	viewOrderCmd = &cobra.Command{
		Use:   "order",
		Short: "View orders",
		Long:  "View orders",
		Run:   viewOrdersCmdRun,
	}
)

func resetViewOrderFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().Uint64VarP(&userID, "userID", "u", 0, "userID (required)")
	cmd.PersistentFlags().Uint64VarP(&orderID, "orderID", "o", 0, "first orderID which should be output")
	cmd.PersistentFlags().Uint64VarP(&ordersLimit, "n", "n", 0, "limit of returned orders")
	cmd.MarkPersistentFlagRequired("userID")
}

func resetViewRefundFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().Uint64VarP(&pageID, "pageID", "p", 0, "pageID (starts from 0)")
	cmd.PersistentFlags().Uint64VarP(&ordersPerPage, "ordersPerPage", "c", 10, "orders per page")
}

func viewRefundCmdRun(cmd *cobra.Command, args []string) {
	defer resetViewRefundFlags(cmd)
	refunds, err := st.GetRefunds(pageID, ordersPerPage)
	if err != nil {
		fmt.Printf("error while view refund: %s\n", err)
		return
	}

	if len(refunds) == 0 {
		fmt.Printf("there are no refunds for page ")
		fmt.Printf("%d with ordersPerPage equal to %d\n", pageID, ordersPerPage)
		return
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderID | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderID }}
{{ "UserID:" | faint }}  {{ .UserID }}
{{ "Date of refund: " | faint }}  {{ .ExpirationDate }}`,
	}

	promt := promptui.Select{
		Label:     "OrderID:",
		Items:     refunds,
		Templates: templates,
	}

	_, _, err = promt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func viewOrdersCmdRun(cmd *cobra.Command, args []string) {
	defer resetViewOrderFlags(cmd)
	orders, err := st.GetOrdersByUserID(userID, orderID, ordersLimit)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(orders) == 0 {
		fmt.Printf("User %d doesn't have orders\n", userID)
		return
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderID | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderID }}
{{ "Expiration date: " | faint }}  {{ .ExpirationDate }}`,
	}

	promt := promptui.Select{
		Label:     fmt.Sprintf("Orders of User %d", userID),
		Items:     orders,
		Templates: templates,
	}

	_, _, err = promt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
