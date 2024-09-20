package storage

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type User struct {
	Orders map[uint64]*domain.Order `json:"orders"`

	OrdersArray     []domain.OrderView `json:"ordersArray"`
	OrdersIDatArray map[uint64]int     `json:"ordersIDatArray"`

	UserID uint64 `json:"userID"`
}

func NewUser(userID uint64) *User {
	return &User{
		Orders:          make(map[uint64]*domain.Order),
		OrdersArray:     make([]domain.OrderView, 0),
		OrdersIDatArray: make(map[uint64]int),
		UserID:          userID,
	}
}

func (u *User) Add(orderID uint64, order *domain.Order) error {
	if _, ok := u.Orders[orderID]; ok {
		return fmt.Errorf("order %d has already accepted", orderID)
	}

	u.Orders[orderID] = order
	u.OrdersArray = append(u.OrdersArray, domain.OrderView{
		Order:   order,
		UserID:  u.UserID,
		OrderID: orderID,
		Exist:   true,
	})
	u.OrdersIDatArray[orderID] = len(u.OrdersArray) - 1

	return nil
}

func (u *User) Get(orderID uint64) (*domain.Order, error) {
	order, ok := u.Orders[orderID]
	if !ok {
		return nil, fmt.Errorf("not found order %d", orderID)
	}

	return order, nil
}

func (u *User) CanRemove(orderID uint64) error {
	_, ok := u.OrdersIDatArray[orderID]
	if !ok {
		return fmt.Errorf("not found order %d at orders array of user %d", orderID, u.UserID)
	}

	return nil
}

func (u *User) Remove(orderId uint64) error {
	id, ok := u.OrdersIDatArray[orderId]
	if !ok {
		return fmt.Errorf("not found order %d at orders array of user %d", orderId, u.UserID)
	}

	u.OrdersArray[id].Exist = false
	delete(u.Orders, orderId)

	return nil
}

func (u *User) GetExpirationDate(orderID uint64) (time.Time, error) {
	order, ok := u.Orders[orderID]
	if !ok {
		return time.Time{}, fmt.Errorf("user %d doesn't have order %d", u.UserID, orderID)
	}

	expDate, err := time.Parse("02-01-2006", order.ExpirationDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing Expiration Date: %w", err)
	}

	return expDate.Truncate(24 * time.Hour).UTC(), nil
}

func (u *User) findID(firstOrderID uint64) (int, error) {
	var ok bool

	id := 0
	if firstOrderID != 0 {
		id, ok = u.OrdersIDatArray[firstOrderID]
		if !ok {
			return 0, fmt.Errorf("not found order %d", firstOrderID)
		}
	}

	return id, nil
}

func (u *User) GetOrders(firstOrderID, limit uint64) ([]domain.OrderView, error) {
	id, err := u.findID(firstOrderID)
	if err != nil {
		return nil, err
	}

	limit = min(uint64(len(u.OrdersArray)), max(limit, uint64(len(u.OrdersArray))))
	res := make([]domain.OrderView, 0)
	orderCount := uint64(0)

	for ; id < int(limit) && orderCount < limit; id++ {
		if u.OrdersArray[id].Exist {
			res = append(res, u.OrdersArray[id])
			orderCount++
		}
	}

	return res, nil
}
