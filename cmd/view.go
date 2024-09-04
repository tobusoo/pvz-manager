package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
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

type refundOrderView struct {
	OrderID, UserID uint64
	Date            string
}

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

	orders := make([]refundOrderView, 0)
	for _, orderID := range refunds {
		stat, err := st.GetOrderStatus(orderID)
		if err != nil {
			continue
		}
		orders = append(orders, refundOrderView{
			OrderID: orderID,
			UserID:  stat.UserID,
			Date:    stat.Date,
		})
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderID | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderID }}
{{ "UserID:" | faint }}  {{ .UserID }}
{{ "Date of refund: " | faint }}  {{ .Date }}`,
	}

	promt := promptui.Select{
		Label:     "OrderID:",
		Items:     orders,
		Templates: templates,
	}

	_, _, err = promt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
