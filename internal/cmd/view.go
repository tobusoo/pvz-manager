package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
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

	viewRefundCmd = &cobra.Command{
		Use:   "refund",
		Short: "View refunds",
		Long:  "View refunds",
		Run:   viewRefundCmdRun,
	}

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
	cmd.PersistentFlags().Uint64VarP(&ordersLimit, "n", "n", 25, "limit of returned orders (defalut 25)")
	cmd.MarkPersistentFlagRequired("userID")
}

func resetViewRefundFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().Uint64VarP(&pageID, "pageID", "p", 1, "pageID (starts from 1)")
	cmd.PersistentFlags().Uint64VarP(&ordersPerPage, "ordersPerPage", "c", 10, "orders per page")
}

func viewRefundCmdRun(cmd *cobra.Command, args []string) {
	defer resetViewRefundFlags(cmd)
	req := &dto.ViewRefundsRequest{
		PageID:        pageID,
		OrdersPerPage: ordersPerPage,
	}

	refunds, err := viewUsecase.GetRefunds(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderID | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderID }} {{ "UserID:" | faint }}  {{ .UserID }}
{{"Cost:" | faint }} {{ .Cost }}rub {{"Weight:" | faint }} {{ .Weight }}gr
{{ "Date of refund: " | faint }}  {{ .ExpirationDate }} {{ "Package Type:" | faint }} {{ .PackageType }}`,
	}

	promt := promptui.Select{
		Label:     "OrderID:",
		Items:     refunds,
		Templates: templates,
	}

	InOutLock()
	_, _, err = promt.Run()
	if err != nil {
		fmt.Println(err)
	}
	InOutUnlock()
}

func viewOrdersCmdRun(cmd *cobra.Command, args []string) {
	defer resetViewOrderFlags(cmd)

	req := &dto.ViewOrdersRequest{
		UserID:       userID,
		FirstOrderID: orderID,
		OrdersLimit:  ordersLimit,
	}

	orders, err := viewUsecase.GetOrders(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderID | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderID }} {{"Cost:" | faint }} {{ .Cost }}rub {{"Weight:" | faint }} {{ .Weight }}gr
{{ "Expiration date:" | faint }} {{ .ExpirationDate }} {{ "Package Type:" | faint }} {{ .PackageType }}`,
	}

	promt := promptui.Select{
		Label:     fmt.Sprintf("Orders of User %d", userID),
		Items:     orders,
		Templates: templates,
	}

	InOutLock()
	_, _, err = promt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	InOutUnlock()
}
