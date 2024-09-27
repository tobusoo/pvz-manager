package storage

import (
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
		GetOrderOnlyStatus(orderID uint64) (stat string, err error)
		SetOrderStatus(orderID uint64, status string) error
	}

	UsersRepository interface {
		AddOrder(userID, orderID uint64, order *domain.Order) error
		GetOrder(userID, orderID uint64) (*domain.Order, error)
		RemoveOrder(orderID uint64, status string) error
		RemoveOrders(ordersID []uint64, status string) error
		CanRemoveOrder(orderID uint64) error
		GetExpirationDate(userID, orderID uint64) (time.Time, error)
		GetOrdersByUserID(userID, firstOrderID, limit uint64) ([]domain.OrderView, error)
	}

	Storage interface {
		RefundsRepository
		OrdersHistoryRepository
		UsersRepository
	}
)
