package domain

import "gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"

const (
	StatusAccepted    = "accepted"
	StatusGiveClient  = "issued to client"
	StatusGiveCourier = "issued to courier"
	StatusReturned    = "returned"
)

type (
	Order struct {
		ExpirationDate string `json:"expirationDate"`
		PackageType    string `json:"packageType"`
		Cost           uint64 `json:"cost"`
		Weight         uint64 `json:"weight"`
		UseTape        bool   `json:"useTape"`
	}

	OrderStatus struct {
		*Order
		Status string `json:"status"`
		Date   string `json:"date"`
		UserID uint64 `json:"userID"`
	}

	OrderView struct {
		*Order
		UserID  uint64 `json:"userID"`
		OrderID uint64 `json:"orderID"`
		Exist   bool   `json:"exist"`
	}
)

func NewOrder(
	cost, weight uint64,
	expDate string,
	cs strategy.ContainerStrategy,
) (order *Order, err error) {

	order = &Order{
		ExpirationDate: expDate,
		Weight:         weight,
		PackageType:    cs.Type(),
		UseTape:        cs.IsTaped(),
	}

	order.Cost, err = cs.CalculateCost(weight, cost)
	return
}
