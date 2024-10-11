package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	manager_service "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
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
	req := &manager_service.ViewRefundsRequestV1{
		PageId:        pageID,
		OrdersPerPage: ordersPerPage,
	}

	refunds, err := mng_service.ViewRefundsV1(ctx, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderId | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderId }} {{ "UserID:" | faint }}  {{ .UserId }}
{{"Cost:" | faint }} {{ .Order.Cost }}rub {{"Weight:" | faint }} {{ .Order.Weight }}gr
{{ "Date of refund: " | faint }}  {{ .Order.ExpirationDate }} {{ "Package Type:" | faint }} {{ .Order.PackageType }}`,
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

	req := &manager_service.ViewOrdersRequestV1{
		UserId:       userID,
		FirstOrderId: orderID,
		Limit:        ordersLimit,
	}

	orders, err := mng_service.ViewOrdersV1(ctx, req)
	if err != nil {
		fmt.Println(err)
		return
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderId | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderId }} {{"Cost:" | faint }} {{ .Order.Cost }}rub {{"Weight:" | faint }} {{ .Order.Weight }}gr
{{ "Expiration date:" | faint }} {{ .Order.ExpirationDate }} {{ "Package Type:" | faint }} {{ .Order.PackageType }}`,
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
