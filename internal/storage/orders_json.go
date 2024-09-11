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

type OrderView struct {
	Order
	UserID  uint64 `json:"userID"`
	OrderID uint64 `json:"orderID"`
	Exist   bool   `json:"exist"`
}

type User struct {
	Orders map[uint64]Order `json:"orders"`

	OrdersArray     []OrderView    `json:"ordersArray"`
	OrdersIDatArray map[uint64]int `json:"ordersIDatArray"`
}

type OrderStatus struct {
	UserID uint64 `json:"userID"`
	Status string `json:"status"`
	Date   string `json:"date"`
}

type Refunds struct {
	Orders          []OrderView    `json:"orders"`
	OrdersIDatArray map[uint64]int `json:"ordersIDatArray"`
}

type Storage struct {
	Users         map[uint64]*User       `json:"users"`
	OrdersHistory map[uint64]OrderStatus `json:"ordersHistory"`
	Refunds       Refunds                `json:"refunds"`

	path string `json:"-"`
}

func NewStorage(path string) (*Storage, error) {
	storage := &Storage{
		path:          path,
		Users:         make(map[uint64]*User),
		OrdersHistory: make(map[uint64]OrderStatus),
		Refunds:       Refunds{make([]OrderView, 0), make(map[uint64]int)},
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

func (s *Storage) SetOrderStatus(orderID uint64, status string) error {
	order, ok := s.OrdersHistory[orderID]
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}

	s.OrdersHistory[orderID] = OrderStatus{order.UserID, status, utils.CurrentDateString()}
	return nil
}

func (s *Storage) GetOrderStatus(orderID uint64) (OrderStatus, error) {
	status, ok := s.OrdersHistory[orderID]
	if !ok {
		return OrderStatus{}, fmt.Errorf("order %d not found", orderID)
	}

	return status, nil
}

func (s *Storage) GetExpirationDate(userID, orderID uint64) (time.Time, error) {
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

	return expDate.Truncate(24 * time.Hour).UTC(), nil
}

func (s *Storage) getRefundsCheckErr(pageID, ordersPerPage uint64) error {
	if ordersPerPage == 0 {
		return fmt.Errorf("orders per page must be greater than 0")
	}

	if pageID == 0 {
		return fmt.Errorf("pageID myst be greater than 0")
	}

	return nil
}

func (s *Storage) getRefundsSlice(firstOrderID int, ordersPerPage uint64) []OrderView {
	res := make([]OrderView, 0)
	ordersCount := 0

	for i := firstOrderID; i < len(s.Refunds.Orders); i++ {
		if s.Refunds.Orders[i].Exist {
			res = append(res, s.Refunds.Orders[i])
			ordersCount++
		}

		if ordersCount == int(ordersPerPage) {
			break
		}
	}

	return res
}

func (s *Storage) GetRefunds(pageID, ordersPerPage uint64) (res []OrderView, err error) {
	if err := s.getRefundsCheckErr(pageID, ordersPerPage); err != nil {
		return nil, err
	}

	firstOrderID := int((pageID - 1) * ordersPerPage)
	return s.getRefundsSlice(firstOrderID, ordersPerPage), nil
}

func (s *Storage) getUserAndFirstOrderID(userID, firstOrderID uint64) (*User, int, error) {
	user, ok := s.Users[userID]
	if !ok {
		return nil, 0, fmt.Errorf("not found user %d", userID)
	}

	id := 0
	if firstOrderID != 0 {
		id, ok = user.OrdersIDatArray[firstOrderID]
		if !ok {
			return nil, 0, fmt.Errorf("not found order %d", firstOrderID)
		}
	}

	return user, id, nil
}

func (s *Storage) GetOrdersByUserID(userID, firstOrderID, limit uint64) ([]OrderView, error) {
	user, id, err := s.getUserAndFirstOrderID(userID, firstOrderID)
	if err != nil {
		return nil, err
	}

	limit = min(uint64(len(user.OrdersArray)), limit)
	res := make([]OrderView, 0)
	orderCount := uint64(0)

	for ; id < int(limit) && orderCount < limit; id++ {
		if user.OrdersArray[id].Exist {
			res = append(res, user.OrdersArray[id])
			orderCount++
		}
	}

	return res, nil
}

func (s *Storage) AddOrder(userID, orderID uint64, expirationDate string) error {
	if order, ok := s.OrdersHistory[orderID]; ok {
		return fmt.Errorf("order %d has already been %s", orderID, order.Status)
	}

	if _, ok := s.Users[userID]; !ok {
		s.Users[userID] = &User{
			make(map[uint64]Order),
			make([]OrderView, 0),
			make(map[uint64]int),
		}
	}

	order := Order{expirationDate}
	s.Users[userID].Orders[orderID] = order
	s.Users[userID].OrdersArray = append(s.Users[userID].OrdersArray, OrderView{order, userID, orderID, true})
	s.Users[userID].OrdersIDatArray[orderID] = len(s.Users[userID].OrdersArray) - 1
	s.OrdersHistory[orderID] = OrderStatus{userID, StatusAccepted, utils.CurrentDateString()}

	return nil
}

func (s *Storage) AddRefund(orderID uint64) (err error) {
	if err = s.SetOrderStatus(orderID, StatusReturned); err != nil {
		return err
	}

	stat, err := s.GetOrderStatus(orderID)
	if err != nil {
		return err
	}

	s.Refunds.Orders = append(s.Refunds.Orders, OrderView{Order{stat.Date}, stat.UserID, orderID, true})
	s.Refunds.OrdersIDatArray[orderID] = len(s.Refunds.Orders) - 1
	return nil
}

func (s *Storage) RemoveReturned(orderID uint64) error {
	id, ok := s.Refunds.OrdersIDatArray[orderID]
	if !ok {
		return fmt.Errorf("not found order %d at refunded", orderID)
	}

	s.Refunds.Orders[id].Exist = false

	return nil
}

func canRemoveOrderCheckStatus(status string, orderID uint64) error {
	if status == StatusGiveClient || status == StatusGiveCourier {
		return fmt.Errorf("order %d has already been %s", orderID, StatusGiveClient)
	}

	return nil
}

func (s *Storage) CanRemoveOrder(orderID uint64) error {
	order, ok := s.OrdersHistory[orderID]
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}

	if err := canRemoveOrderCheckStatus(order.Status, orderID); err != nil {
		return err
	}

	user, ok := s.Users[order.UserID]
	if !ok {
		return fmt.Errorf("user %d not found", order.UserID)
	}

	_, ok = user.OrdersIDatArray[orderID]
	if !ok {
		return fmt.Errorf("not found order %d at orders array of user %d", orderID, order.UserID)
	}

	return nil
}

// Использовать только перед вызовом CanRemoveOrder!!!
func (s *Storage) RemoveOrder(orderID uint64, status string) error {
	order := s.OrdersHistory[orderID]
	user := s.Users[order.UserID]
	id := user.OrdersIDatArray[orderID]

	s.OrdersHistory[orderID] = OrderStatus{order.UserID, status, utils.CurrentDateString()}
	user.OrdersArray[id].Exist = false
	delete(user.Orders, orderID)
	return nil
}
