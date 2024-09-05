package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	viewCmd.AddCommand(viewRefundCmd)
	viewCmd.AddCommand(viewOrderCmd)

	resetViewOrderFlags(viewOrderCmd)
	viewOrderCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetViewOrderFlags(cmd)
	})

	viewCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		cmd.ResetFlags()
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

	ordersCount  uint
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
	cmd.PersistentFlags().UintVarP(&ordersCount, "n", "n", 0, "print maximum N last orders")
	cmd.MarkPersistentFlagRequired("userID")
}

type orderView struct {
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

	orders := make([]orderView, 0)
	for _, orderID := range refunds {
		stat, err := st.GetOrderStatus(orderID)
		if err != nil {
			continue
		}
		orders = append(orders, orderView{
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

func viewOrdersCmdRun(cmd *cobra.Command, args []string) {
	defer resetViewOrderFlags(cmd)
	ordersMap, err := st.GetOrdersByUserID(userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(ordersMap) == 0 {
		fmt.Printf("User %d doesn't have orders\n", userID)
		return
	}

	orders := make([]orderView, 0)
	for orderID, v := range ordersMap {
		orders = append(orders, orderView{
			OrderID: orderID,
			UserID:  userID,
			Date:    v.ExpirationDate,
		})
	}

	if ordersCount > 0 && int(ordersCount) < len(orders) {
		orders = sortOrdersByDate(orders)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{.}}",
		Active:   "\U0001F336 {{.OrderID | cyan}}",
		Inactive: "  {{.OrderID | cyan}}",
		Selected: " ",
		Details: `-----Order-----
{{ "OrderID:" | faint }}  {{ .OrderID }}
{{ "Expiration date: " | faint }}  {{ .Date }}`,
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

func sortOrdersByDate(orders []orderView) []orderView {
	sort.Slice(orders, func(i, j int) bool {
		timeStr1 := orders[i].Date
		timeStr2 := orders[j].Date
		time1, err := time.Parse("02-01-2006", timeStr1)
		if err != nil {
			// Bad data at storage
			fmt.Printf("Error while parsing time: %s\n", err)
			os.Exit(1)
		}
		time2, err := time.Parse("02-01-2006", timeStr2)
		if err != nil {
			fmt.Printf("Error while parsing time: %s\n", err)
			os.Exit(1)
		}

		return time1.Before(time2)
	})

	return orders[len(orders)-int(ordersCount):]
}
