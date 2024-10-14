package domain

import (
	"errors"

	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
)

const (
	StatusAccepted    = "accepted"
	StatusGiveClient  = "issued to client"
	StatusGiveCourier = "issued to courier"
	StatusReturned    = "returned"
)

var (
	ErrNotFound             = errors.New("order not found")
	ErrWrongInput           = errors.New("wrong input")
	ErrWrongStatus          = errors.New("wrong status")
	ErrAlreadyExist         = errors.New("order already exist")
	ErrNotExpirationDate    = errors.New("expiration date hasn't expired yet")
	ErrExpirationDatePassed = errors.New("expiration date has already passed")
	ErrTwoDaysPassed        = errors.New("2 days have passed since the order was issued to the client")
)

type (
	Order struct {
		ExpirationDate string `json:"expirationDate" db:"expiration_date"`
		PackageType    string `json:"packageType" db:"package_type"`
		Cost           uint64 `json:"cost" db:"cost"`
		Weight         uint64 `json:"weight" db:"weight"`
		UseTape        bool   `json:"useTape" db:"use_tape"`
	}

	OrderStatus struct {
		*Order
		Status    string `json:"status" db:"status"`
		UpdatedAt string `json:"updatedAt" db:"updated_at"`
		UserID    uint64 `json:"userID" db:"user_id"`
	}

	OrderView struct {
		*Order
		UserID  uint64 `json:"userID" db:"user_id"`
		OrderID uint64 `json:"orderID" db:"order_id"`
		Exist   bool   `json:"exist" db:"-"`
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
