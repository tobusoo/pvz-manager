package storage_json

import (
	"fmt"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

type OrdersHistory struct {
	Stat map[uint64]*domain.OrderStatus `json:"ordersHistory"`
}

func NewOrdersHistory() *OrdersHistory {
	return &OrdersHistory{make(map[uint64]*domain.OrderStatus)}
}

func (s *OrdersHistory) AddOrderStatus(orderID, userID uint64, status string, order *domain.Order) error {
	stat, ok := s.Stat[orderID]
	if ok {
		return fmt.Errorf("order %d has already been %s", orderID, stat.Status)
	}

	s.Stat[orderID] = &domain.OrderStatus{
		Order:  order,
		Status: status,
		Date:   utils.CurrentDateString(),
		UserID: userID,
	}

	return nil
}

func (s *OrdersHistory) GetOrderStatus(orderID uint64) (*domain.OrderStatus, error) {
	status, ok := s.Stat[orderID]
	if !ok {
		return nil, fmt.Errorf("order %d not found", orderID)
	}

	return status, nil
}

func (s *OrdersHistory) SetOrderStatus(orderID uint64, status string) error {
	order, ok := s.Stat[orderID]
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}

	order.Status = status
	return nil
}
