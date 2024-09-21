package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type (
	RefundsRepository interface {
		AddRefund(userID, orderID uint64, order *domain.Order) error
		RemoveRefund(orderID uint64) error
		GetRefunds(pageID, ordersPerPage uint64) ([]domain.OrderView, error)
	}

	OrdersHistoryRepository interface {
		AddOrderStatus(orderID, userID uint64, status string, order *domain.Order) error
		GetOrderStatus(orderID uint64) (*domain.OrderStatus, error)
		SetOrderStatus(orderID uint64, status string) error
	}

	UsersRepository interface {
		AddOrder(userID, orderID uint64, order *domain.Order) error
		GetOrder(userID, orderID uint64) (*domain.Order, error)
		CanRemove(userID, orderID uint64) error
		RemoveOrder(userID, orderID uint64) error
		GetExpirationDate(userID, orderID uint64) (time.Time, error)
		GetOrders(userID, firstOrderID, limit uint64) ([]domain.OrderView, error)
	}

	Storage interface {
		RefundsRepository
		OrdersHistoryRepository
		AddOrder(userID, orderID uint64, order *domain.Order) error
		GetOrder(userID, orderID uint64) (*domain.Order, error)
		GetExpirationDate(userID, orderID uint64) (time.Time, error)
		GetOrdersByUserID(userID, firstOrderID, limit uint64) ([]domain.OrderView, error)
		CanRemoveOrder(orderID uint64) error
		RemoveOrder(orderID uint64, status string) error
	}
)

type StorageJSON struct {
	OrdersHistoryRepository `json:"historyRepository"`
	RefundsRepository       `json:"refundsRepository"`
	Users                   UsersRepository `json:"usersRepository"`

	path string `json:"-"`
}

func NewStorage(
	ohp OrdersHistoryRepository,
	rp RefundsRepository,
	up UsersRepository,
	path string,
) (*StorageJSON, error) {
	storage := &StorageJSON{
		OrdersHistoryRepository: ohp,
		RefundsRepository:       rp,
		Users:                   up,
		path:                    path,
	}

	err := storage.readDataFromFile()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *StorageJSON) readDataFromFile() (err error) {
	file, err := os.OpenFile(s.path, os.O_RDWR, 0666)
	if err != nil {
		file, err = os.Create(s.path)
		file.WriteString("{}")
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	return
}

func (s *StorageJSON) Save() (err error) {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s)
	return
}

func (s *StorageJSON) AddOrder(userID, orderID uint64, order *domain.Order) error {
	stat, err := s.GetOrderStatus(orderID)
	if err == nil {
		return fmt.Errorf("order %d has already been %s", orderID, stat.Status)
	}

	err = s.Users.AddOrder(userID, orderID, order)
	if err != nil {
		return err
	}

	return s.AddOrderStatus(orderID, userID, domain.StatusAccepted, order)
}

func (s *StorageJSON) GetOrder(userID, orderID uint64) (*domain.Order, error) {
	return s.Users.GetOrder(userID, orderID)
}

func (s *StorageJSON) GetExpirationDate(userID, orderID uint64) (time.Time, error) {
	return s.Users.GetExpirationDate(userID, orderID)
}

func (s *StorageJSON) GetOrdersByUserID(userID, firstOrderID, limit uint64) ([]domain.OrderView, error) {
	return s.Users.GetOrders(userID, firstOrderID, limit)
}

func canRemoveOrderCheckStatus(status string, orderID uint64) error {
	if status == domain.StatusGiveClient || status == domain.StatusGiveCourier {
		return fmt.Errorf("order %d has already been %s", orderID, domain.StatusGiveClient)
	}

	return nil
}

func (s *StorageJSON) CanRemoveOrder(orderID uint64) error {
	stat, err := s.GetOrderStatus(orderID)
	if err != nil {
		return err
	}

	if err = canRemoveOrderCheckStatus(stat.Status, orderID); err != nil {
		return err
	}

	return s.Users.CanRemove(stat.UserID, orderID)
}

// Использовать только перед вызовом CanRemoveOrder!!!
func (s *StorageJSON) RemoveOrder(orderID uint64, status string) error {
	stat, _ := s.GetOrderStatus(orderID)

	err := s.Users.RemoveOrder(stat.UserID, orderID)
	if err != nil {
		return err
	}

	s.SetOrderStatus(orderID, status)
	return nil
}
