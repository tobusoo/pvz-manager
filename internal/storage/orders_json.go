package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

const (
	StatusAccepted    = "accepted"
	StatusGiveClient  = "issued to client"
	StatusGiveCourier = "issued to courier"
	StatusReturned    = "returned"
)

type Order struct {
	ExpirationDate string `json:"expirationDate"`
}

type User struct {
	Orders map[uint64]Order `json:"orders"`
}

type OrderStatus struct {
	UserID uint64 `json:"userID"`
	Status string `json:"status"`
	Date   string `json:"date"`
}

type Storage struct {
	Users         map[uint64]User        `json:"users"`
	OrdersHistory map[uint64]OrderStatus `json:"ordersHistory"`
	Refunds       map[uint64]struct{}    `json:"refunds"`

	path string `json:"-"`
}

func NewStorage(path string) *Storage {
	return &Storage{
		path:          path,
		Users:         make(map[uint64]User),
		OrdersHistory: make(map[uint64]OrderStatus),
		Refunds:       make(map[uint64]struct{}),
	}
}

func (s *Storage) readDataFromFile() (err error) {
	clear(s.Users)
	clear(s.OrdersHistory)
	clear(s.Refunds)

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

func (s *Storage) writeDataToFile() (err error) {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s)
	return
}

func (s *Storage) SetOrderStatus(orderID uint64, status string) error {
	if err := s.readDataFromFile(); err != nil {
		return err
	}

	order, ok := s.OrdersHistory[orderID]
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}

	s.OrdersHistory[orderID] = OrderStatus{order.UserID, status, utils.CurrentDateString()}
	return s.writeDataToFile()
}

func (s *Storage) GetOrderStatus(orderID uint64) (OrderStatus, error) {
	if err := s.readDataFromFile(); err != nil {
		return OrderStatus{}, err
	}

	status, ok := s.OrdersHistory[orderID]
	if !ok {
		return OrderStatus{}, fmt.Errorf("order %d not found", orderID)
	}

	return status, nil
}

func (s *Storage) GetExpirationDate(userID, orderID uint64) (time.Time, error) {
	if err := s.readDataFromFile(); err != nil {
		return time.Time{}, err
	}

	user, ok := s.Users[userID]
	if !ok {
		return time.Time{}, fmt.Errorf("user %d not found", userID)
	}

	order, ok := user.Orders[orderID]
	if !ok {
		return time.Time{}, fmt.Errorf("user %d doesn't have order %d", userID, orderID)
	}

	expDate, err := time.Parse("02-01-2006", order.ExpirationDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing Expiration Date: %w", err)
	}

	return expDate.Truncate(24 * time.Hour), s.writeDataToFile()
}

func (s *Storage) AddOrder(userID, orderID uint64, expirationDate string) error {
	if err := s.readDataFromFile(); err != nil {
		return err
	}

	if order, ok := s.OrdersHistory[orderID]; ok {
		return fmt.Errorf("order %d has already been %s", orderID, order.Status)
	}

	if _, ok := s.Users[userID]; !ok {
		s.Users[userID] = User{make(map[uint64]Order)}
	}

	s.Users[userID].Orders[orderID] = Order{expirationDate}
	s.OrdersHistory[orderID] = OrderStatus{userID, StatusAccepted, utils.CurrentDateString()}

	return s.writeDataToFile()
}

func (s *Storage) AddRefund(orderID uint64) (err error) {
	if err = s.readDataFromFile(); err != nil {
		return err
	}

	if err = s.SetOrderStatus(orderID, StatusReturned); err != nil {
		return err
	}

	s.Refunds[orderID] = struct{}{}
	return s.writeDataToFile()
}

func (s *Storage) RemoveReturned(orderID uint64) error {
	if err := s.readDataFromFile(); err != nil {
		return err
	}
	delete(s.Refunds, orderID)

	return s.writeDataToFile()
}

func (s *Storage) RemoveOrder(orderID uint64, status string) error {
	if err := s.readDataFromFile(); err != nil {
		return err
	}

	order, ok := s.OrdersHistory[orderID]
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}

	if order.Status == StatusGiveClient || order.Status == StatusGiveCourier {
		return fmt.Errorf("order %d has already been %s", orderID, StatusGiveClient)
	}

	user, ok := s.Users[order.UserID]
	if !ok {
		return fmt.Errorf("user %d not found", order.UserID)
	}

	s.OrdersHistory[orderID] = OrderStatus{order.UserID, status, utils.CurrentDateString()}

	delete(user.Orders, orderID)
	return s.writeDataToFile()
}
