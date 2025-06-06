package storage_json

import (
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type User struct {
	Orders map[uint64]*domain.Order `json:"orders"`

	OrdersArray     []domain.OrderView `json:"ordersArray"`
	OrdersIDatArray map[uint64]int     `json:"ordersIDatArray"`

	UserID uint64 `json:"userID"`

	mtx sync.Mutex
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
	u.mtx.Lock()
	defer u.mtx.Unlock()

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
	u.mtx.Lock()
	order, ok := u.Orders[orderID]
	u.mtx.Unlock()
	if !ok {
		return nil, fmt.Errorf("not found order %d", orderID)
	}

	return order, nil
}

func (u *User) CanRemove(orderID uint64) error {
	u.mtx.Lock()
	_, ok := u.OrdersIDatArray[orderID]
	u.mtx.Unlock()
	if !ok {
		return fmt.Errorf("not found order %d at orders array of user %d", orderID, u.UserID)
	}

	return nil
}

func (u *User) Remove(orderId uint64) error {
	u.mtx.Lock()
	id, ok := u.OrdersIDatArray[orderId]
	if !ok {
		return fmt.Errorf("not found order %d at orders array of user %d", orderId, u.UserID)
	}

	u.OrdersArray[id].Exist = false
	delete(u.Orders, orderId)
	u.mtx.Unlock()

	return nil
}

func (u *User) GetExpirationDate(orderID uint64) (time.Time, error) {
	u.mtx.Lock()
	order, ok := u.Orders[orderID]
	u.mtx.Unlock()
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
		u.mtx.Lock()
		id, ok = u.OrdersIDatArray[firstOrderID]
		u.mtx.Unlock()
		if !ok {
			return 0, fmt.Errorf("not found order %d", firstOrderID)
		}
	}

	return id, nil
}

func (u *User) calcLimit(limit, arrayLen uint64) uint64 {
	if limit == 0 {
		return arrayLen
	}

	return min(limit, arrayLen)
}

func (u *User) GetOrders(firstOrderID, limit uint64) ([]domain.OrderView, error) {
	id, err := u.findID(firstOrderID)
	if err != nil {
		return nil, err
	}

	limit = u.calcLimit(limit, uint64(len(u.OrdersArray)))
	res := make([]domain.OrderView, 0)
	orderCount := uint64(0)

	for ; id < len(u.OrdersArray) && orderCount < limit; id++ {
		if u.OrdersArray[id].Exist {
			res = append(res, u.OrdersArray[id])
			orderCount++
		}
	}

	return res, nil
}
