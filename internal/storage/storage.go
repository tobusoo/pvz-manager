package storage

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
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
		GetExpirationDate(userID, orderID uint64) (time.Time, error)
		GetOrdersByUserID(userID, firstOrderID, limit uint64) ([]domain.OrderView, error)
		CanRemoveOrder(orderID uint64) error
		RemoveOrder(orderID uint64, status string) error
		RemoveOrders(ordersID []uint64, status string) error
	}

	Storage interface {
		RefundsRepository
		OrdersHistoryRepository
		UsersRepository
	}

	StorageDB struct {
		txManager postgres.TransactionManager
		db        Storage
		ctx       context.Context
	}
)

func NewStorageDB(tx postgres.TransactionManager, db Storage) *StorageDB {
	return &StorageDB{
		txManager: tx,
		db:        db,
		ctx:       context.Background(),
	}
}

func (s *StorageDB) AddOrder(userID, orderID uint64, order *domain.Order) (err error) {
	if stat, err := s.GetOrderStatus(orderID); err == nil {
		return fmt.Errorf("order %d has already been %s", orderID, stat.Status)
	}

	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		if err = s.db.AddOrder(userID, orderID, order); err != nil {
			return err
		}
		return s.db.AddOrderStatus(orderID, userID, domain.StatusAccepted, order)
	})
}

func (s *StorageDB) GetOrder(userID, orderID uint64) (order *domain.Order, err error) {
	s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		order, err = s.db.GetOrder(userID, orderID)
		return err
	})
	return
}

func (s *StorageDB) GetExpirationDate(userID, orderID uint64) (t time.Time, err error) {
	s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		t, err = s.db.GetExpirationDate(userID, orderID)
		return err
	})
	return
}

func (s *StorageDB) GetOrdersByUserID(userID, firstOrderID, limit uint64) (orders []domain.OrderView, err error) {
	s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		orders, err = s.db.GetOrdersByUserID(userID, firstOrderID, limit)
		return err
	})
	return
}

func canRemoveOrderCheckStatus(status string, orderID uint64) error {
	if status == domain.StatusGiveClient || status == domain.StatusGiveCourier {
		return fmt.Errorf("order %d has already been %s", orderID, domain.StatusGiveClient)
	}
	return nil
}

func (s *StorageDB) CanRemoveOrder(orderID uint64) error {
	stat, err := s.GetOrderStatus(orderID)
	if err != nil {
		return err
	}

	if err = canRemoveOrderCheckStatus(stat.Status, orderID); err != nil {
		return err
	}

	return s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		return s.db.CanRemoveOrder(orderID)
	})
}

func (s *StorageDB) RemoveOrder(orderID uint64, status string) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		return s.db.RemoveOrder(orderID, status)
	})
}

func (s *StorageDB) RemoveOrders(ordersID []uint64, status string) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		return s.db.RemoveOrders(ordersID, status)
	})
}

func (s *StorageDB) AddOrderStatus(orderID, userID uint64, status string, order *domain.Order) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		return s.db.AddOrderStatus(orderID, userID, status, order)
	})
}

func (s *StorageDB) GetOrderStatus(orderID uint64) (order *domain.OrderStatus, err error) {
	err = s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		order, err = s.db.GetOrderStatus(orderID)
		return err
	})
	return
}

func (s *StorageDB) SetOrderStatus(orderID uint64, status string) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		return s.db.SetOrderStatus(orderID, status)
	})
}

func (s *StorageDB) AddRefund(userID, orderID uint64, order *domain.Order) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		err := s.db.AddRefund(userID, orderID, order)
		if err != nil {
			return err
		}
		return s.db.SetOrderStatus(orderID, domain.StatusReturned)
	})
}

func (s *StorageDB) RemoveRefund(orderID uint64) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		err := s.db.RemoveRefund(orderID)
		if err != nil {
			return err
		}
		return s.db.SetOrderStatus(orderID, domain.StatusGiveCourier)
	})
}

func (s *StorageDB) GetRefunds(pageID, ordersPerPage uint64) (orders []domain.OrderView, err error) {
	err = s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		orders, err = s.db.GetRefunds(pageID, ordersPerPage)
		if err != nil {
			return err
		}
		return nil
	})
	return orders, err
}
