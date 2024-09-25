package storage_json

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

type UsersRepository interface {
	AddOrder(userID, orderID uint64, order *domain.Order) error
	GetOrder(userID, orderID uint64) (*domain.Order, error)
	CanRemove(userID, orderID uint64) error
	RemoveOrder(userID, orderID uint64) error
	GetExpirationDate(userID, orderID uint64) (time.Time, error)
	GetOrders(userID, firstOrderID, limit uint64) ([]domain.OrderView, error)
}

type Storage struct {
	Ohp   storage.OrdersHistoryRepository `json:"historyRepository"`
	Rp    storage.RefundsRepository       `json:"refundsRepository"`
	Users UsersRepository                 `json:"usersRepository"`

	path string `json:"-"`
}

func NewStorage(
	ohp storage.OrdersHistoryRepository,
	rp storage.RefundsRepository,
	up UsersRepository,
	path string,
) (*Storage, error) {
	storage := &Storage{
		Ohp:   ohp,
		Rp:    rp,
		Users: up,
		path:  path,
	}

	err := storage.readDataFromFile()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) readDataFromFile() (err error) {
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

func (s *Storage) Save() (err error) {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s)
	return
}

func (s *Storage) AddOrderStatus(orderID, userID uint64, status string, order *domain.Order) error {
	return s.Ohp.AddOrderStatus(orderID, userID, status, order)
}

func (s *Storage) GetOrderStatus(orderID uint64) (*domain.OrderStatus, error) {
	return s.Ohp.GetOrderStatus(orderID)
}

func (s *Storage) SetOrderStatus(orderID uint64, status string) error {
	return s.Ohp.SetOrderStatus(orderID, status)
}

func (s *Storage) AddRefund(userID, orderID uint64, order *domain.Order) error {
	if err := s.Rp.AddRefund(userID, orderID, order); err != nil {
		return err
	}

	return s.SetOrderStatus(orderID, domain.StatusReturned)
}

func (s *Storage) RemoveRefund(orderID uint64) error {
	if err := s.Rp.RemoveRefund(orderID); err != nil {
		return err
	}

	return s.SetOrderStatus(orderID, domain.StatusGiveCourier)
}

func (s *Storage) GetRefunds(pageID, ordersPerPage uint64) ([]domain.OrderView, error) {
	return s.Rp.GetRefunds(pageID, ordersPerPage)
}

func (s *Storage) AddOrder(userID, orderID uint64, order *domain.Order) error {
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

func (s *Storage) GetOrder(userID, orderID uint64) (*domain.Order, error) {
	return s.Users.GetOrder(userID, orderID)
}

func (s *Storage) GetExpirationDate(userID, orderID uint64) (time.Time, error) {
	return s.Users.GetExpirationDate(userID, orderID)
}

func (s *Storage) GetOrdersByUserID(userID, firstOrderID, limit uint64) ([]domain.OrderView, error) {
	return s.Users.GetOrders(userID, firstOrderID, limit)
}

func canRemoveOrderCheckStatus(status string, orderID uint64) error {
	if status == domain.StatusGiveClient || status == domain.StatusGiveCourier {
		return fmt.Errorf("order %d has already been %s", orderID, domain.StatusGiveClient)
	}

	return nil
}

func (s *Storage) CanRemoveOrder(orderID uint64) error {
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
func (s *Storage) RemoveOrder(orderID uint64, status string) error {
	stat, _ := s.GetOrderStatus(orderID)

	err := s.Users.RemoveOrder(stat.UserID, orderID)
	if err != nil {
		return err
	}

	s.SetOrderStatus(orderID, status)
	return nil
}

func (s *Storage) RemoveOrders(ordersID []uint64, status string) error {
	for _, orderID := range ordersID {
		s.RemoveOrder(orderID, status)
	}
	return nil
}
