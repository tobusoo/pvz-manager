package storage_json

import (
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type Users struct {
	UsersMap map[uint64]*User `json:"users"`
	mtx      sync.Mutex
}

func NewUsers() *Users {
	return &Users{UsersMap: make(map[uint64]*User)}
}

func (u *Users) AddOrder(userID, orderID uint64, order *domain.Order) error {

	u.mtx.Lock()
	if _, ok := u.UsersMap[userID]; !ok {
		u.UsersMap[userID] = NewUser(userID)
	}
	u.mtx.Unlock()

	return u.UsersMap[userID].Add(orderID, order)
}

func (u *Users) GetOrder(userID, orderID uint64) (*domain.Order, error) {

	u.mtx.Lock()
	user, ok := u.UsersMap[userID]
	u.mtx.Unlock()
	if !ok {
		return nil, fmt.Errorf("user %d not found", userID)
	}

	return user.Get(orderID)
}

func (u *Users) CanRemove(userID, orderID uint64) error {

	u.mtx.Lock()
	user, ok := u.UsersMap[userID]
	u.mtx.Unlock()
	if !ok {
		return fmt.Errorf("user %d not found", userID)
	}

	return user.CanRemove(orderID)
}

func (u *Users) RemoveOrder(userID, orderID uint64) error {
	u.mtx.Lock()
	user := u.UsersMap[userID]
	u.mtx.Unlock()

	return user.Remove(orderID)
}

func (u *Users) GetExpirationDate(userID, orderID uint64) (time.Time, error) {

	u.mtx.Lock()
	user, ok := u.UsersMap[userID]
	u.mtx.Unlock()
	if !ok {
		return time.Time{}, fmt.Errorf("user %d not found", userID)
	}

	return user.GetExpirationDate(orderID)
}

func (u *Users) GetOrders(userID, firstOrderID, limit uint64) ([]domain.OrderView, error) {

	u.mtx.Lock()
	user, ok := u.UsersMap[userID]
	u.mtx.Unlock()
	if !ok {
		return nil, fmt.Errorf("not found user %d", userID)
	}

	return user.GetOrders(firstOrderID, limit)
}
