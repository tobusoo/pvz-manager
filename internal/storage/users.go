package storage

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type Users struct {
	UsersMap map[uint64]*User `json:"users"`
}

func NewUsers() *Users {
	return &Users{make(map[uint64]*User)}
}

func (u *Users) AddOrder(userID, orderID uint64, order *domain.Order) error {
	if _, ok := u.UsersMap[userID]; !ok {
		u.UsersMap[userID] = NewUser(userID)
	}

	return u.UsersMap[userID].Add(orderID, order)
}

func (u *Users) CanRemove(userID, orderID uint64) error {
	user, ok := u.UsersMap[userID]
	if !ok {
		return fmt.Errorf("user %d not found", userID)
	}

	return user.CanRemove(orderID)
}

func (u *Users) RemoveOrder(userID, orderID uint64) error {
	user := u.UsersMap[userID]
	return user.Remove(orderID)
}

func (u *Users) GetExpirationDate(userID, orderID uint64) (time.Time, error) {
	user, ok := u.UsersMap[userID]
	if !ok {
		return time.Time{}, fmt.Errorf("user %d not found", userID)
	}

	return user.GetExpirationDate(orderID)
}

func (u *Users) GetOrders(userID, firstOrderID, limit uint64) ([]domain.OrderView, error) {
	user, ok := u.UsersMap[userID]
	if !ok {
		return nil, fmt.Errorf("not found user %d", userID)
	}

	return user.GetOrders(firstOrderID, limit)
}
