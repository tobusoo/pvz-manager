package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

func init() {
	acceptCmd.AddCommand(acceptOrderCmd)
	acceptCmd.AddCommand(acceptRefundCmd)

	resetOrderFlags(acceptOrderCmd)
	acceptOrderCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetOrderFlags(cmd)
	})

	resetRefundFlags(acceptRefundCmd)
	acceptRefundCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		resetRefundFlags(cmd)
	})
}

var (
	acceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "Accept orders or refund",
		Long:  "Accept orders or refund",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cost, weight   uint64
	containerType  string
	useTape        bool
	acceptOrderCmd = &cobra.Command{
		Use:   "order",
		Short: "Accept order",
		Long:  "Accept order",
		Run:   acceptOrderCmdRun,
	}

	acceptRefundCmd = &cobra.Command{
		Use:   "refund",
		Short: "Accept refund",
		Long:  "Accept refund",
		Run:   acceptRefundCmdRun,
	}
)

func resetRefundFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.PersistentFlags().Uint64VarP(&userID, "userID", "u", 0, "userID (required)")
	cmd.PersistentFlags().Uint64VarP(&orderID, "orderID", "o", 0, "orderID (required)")
	cmd.MarkPersistentFlagRequired("userID")
	cmd.MarkPersistentFlagRequired("orderID")
}

func resetOrderFlags(cmd *cobra.Command) {
	resetRefundFlags(cmd)

	cmd.PersistentFlags().Uint64VarP(&cost, "cost", "c", 0, "cost in rubles (required)")
	cmd.PersistentFlags().Uint64VarP(&weight, "weight", "w", 0, "weight in grams (required)")
	cmd.PersistentFlags().StringVarP(&expirationDate, "time", "t", "", "Expiration Date (required)")
	cmd.MarkPersistentFlagRequired("cost")
	cmd.MarkPersistentFlagRequired("weight")
	cmd.MarkPersistentFlagRequired("time")

	cmd.PersistentFlags().StringVarP(&containerType, "containerType", "p", "", "containerType (tape, package, box)")
	cmd.PersistentFlags().BoolVarP(&useTape, "useTape", "s", false, "use additional tape")
}

func generateOrder(expDate string) (*domain.Order, error) {
	var cs strategy.ContainerStrategy
	cs, ok := strategy.ContainerTypeMap[containerType]
	if !ok {
		return nil, fmt.Errorf("%s isn't container type", containerType)
	}

	if useTape {
		err := cs.UseTape()
		if err != nil {
			return nil, err
		}
	}

	return domain.NewOrder(cost, weight, expDate, cs)
}

func acceptOrderCmdRun(cmd *cobra.Command, args []string) {
	defer resetOrderFlags(cmd)

	expDate, err := time.Parse("02-01-2006", expirationDate)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentDate := utils.CurrentDate()
	if currentDate.After(expDate) {
		fmt.Println("expiration date has already passed")
		return
	}

	order, err := generateOrder(expDate.Format("02-01-2006"))
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = st.AddOrder(userID, orderID, order); err != nil {
		fmt.Println(err)
	}
}

func acceptRefundCheckErr(order *domain.OrderStatus) error {
	if order.Status != domain.StatusGiveClient {
		return fmt.Errorf("can not refund order %d: status = %s", orderID, order.Status)
	}

	if userID != order.UserID {
		return fmt.Errorf("can not refund order %d: wrong userID", orderID)
	}

	issuedDate, err := time.Parse("02-01-2006", order.Date)
	if err != nil {
		return err
	}

	issuedDate = issuedDate.Add(2 * 24 * time.Hour)
	currentDate := utils.CurrentDate()

	if currentDate.After(issuedDate) {
		return fmt.Errorf("2 days have passed since the order was issued to the client")
	}

	return nil
}

func acceptRefundCmdRun(cmd *cobra.Command, args []string) {
	defer resetRefundFlags(cmd)

	order, err := st.GetOrderStatus(orderID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = acceptRefundCheckErr(order); err != nil {
		fmt.Println(err)
		return
	}

	if err = st.AddRefund(order.UserID, orderID, order.Order); err != nil {
		fmt.Println(err)
		return
	}

	if err = st.SetOrderStatus(orderID, domain.StatusReturned); err != nil {
		fmt.Println(err)
	}
}
