package postgres

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type (
	RefundsRepositoryDB interface {
		AddRefund(ctx context.Context, userID, orderID uint64, order *domain.Order) error
		RemoveRefund(ctx context.Context, orderID uint64) error
		GetRefunds(ctx context.Context, pageID, ordersPerPage uint64) ([]domain.OrderView, error)
	}

	OrdersHistoryRepositoryDB interface {
		AddOrderStatus(ctx context.Context, orderID, userID uint64, status string, order *domain.Order) error
		GetOrderStatus(ctx context.Context, orderID uint64) (*domain.OrderStatus, error)
		SetOrderStatus(ctx context.Context, orderID uint64, status string) error
	}

	UsersRepositoryDB interface {
		AddOrder(ctx context.Context, userID, orderID uint64, order *domain.Order) error
		GetOrder(ctx context.Context, userID, orderID uint64) (*domain.Order, error)
		GetExpirationDate(ctx context.Context, userID, orderID uint64) (time.Time, error)
		GetOrdersByUserID(ctx context.Context, userID, firstOrderID, limit uint64) ([]domain.OrderView, error)
		CanRemoveOrder(ctx context.Context, userID, orderID uint64) error
		RemoveOrder(ctx context.Context, userID, orderID uint64) error
	}

	RepositoryDB interface {
		RefundsRepositoryDB
		OrdersHistoryRepositoryDB
		UsersRepositoryDB
	}

	StorageDB struct {
		txManager TransactionManager
		db        RepositoryDB
		ctx       context.Context
	}
)

func NewStorageDB(ctx context.Context, tx TransactionManager, db RepositoryDB) *StorageDB {
	return &StorageDB{
		txManager: tx,
		db:        db,
		ctx:       ctx,
	}
}

func (s *StorageDB) AddOrder(userID, orderID uint64, order *domain.Order) (err error) {
	if stat, err := s.GetOrderStatus(orderID); err == nil {
		return fmt.Errorf("order %d has already been %s", orderID, stat.Status)
	}

	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		if err = s.db.AddOrder(ctxTx, userID, orderID, order); err != nil {
			return err
		}
		return s.db.AddOrderStatus(ctxTx, orderID, userID, domain.StatusAccepted, order)
	})
}

func (s *StorageDB) GetOrder(userID, orderID uint64) (order *domain.Order, err error) {
	s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		order, err = s.db.GetOrder(ctxTx, userID, orderID)
		return err
	})
	return
}

func (s *StorageDB) GetExpirationDate(userID, orderID uint64) (t time.Time, err error) {
	s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		t, err = s.db.GetExpirationDate(ctxTx, userID, orderID)
		return err
	})
	return
}

func (s *StorageDB) GetOrdersByUserID(userID, firstOrderID, limit uint64) (orders []domain.OrderView, err error) {
	s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		orders, err = s.db.GetOrdersByUserID(ctxTx, userID, firstOrderID, limit)
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
		return s.db.CanRemoveOrder(ctxTx, stat.UserID, orderID)
	})
}

func (s *StorageDB) RemoveOrder(orderID uint64, status string) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		stat, err := s.db.GetOrderStatus(s.ctx, orderID)
		if err != nil {
			return err
		}

		return s.db.RemoveOrder(ctxTx, stat.UserID, orderID)
	})
}

func (s *StorageDB) removeOrder(ctxTx context.Context, orderID uint64, status string) error {
	stat, err := s.db.GetOrderStatus(s.ctx, orderID)
	if err != nil {
		return err
	}

	if err = s.db.RemoveOrder(ctxTx, stat.UserID, orderID); err != nil {
		return err
	}

	return s.db.SetOrderStatus(ctxTx, orderID, status)
}

func (s *StorageDB) RemoveOrders(ordersID []uint64, status string) error {
	return s.txManager.RunSerializable(s.ctx, func(ctxTx context.Context) error {
		for _, order := range ordersID {
			if err := s.removeOrder(s.ctx, order, status); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *StorageDB) AddOrderStatus(orderID, userID uint64, status string, order *domain.Order) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		return s.db.AddOrderStatus(ctxTx, orderID, userID, status, order)
	})
}

func (s *StorageDB) GetOrderStatus(orderID uint64) (order *domain.OrderStatus, err error) {
	err = s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		order, err = s.db.GetOrderStatus(ctxTx, orderID)
		return err
	})
	return
}

func (s *StorageDB) SetOrderStatus(orderID uint64, status string) error {
	return s.txManager.RunSerializable(s.ctx, func(ctxTx context.Context) error {
		return s.db.SetOrderStatus(ctxTx, orderID, status)
	})
}

func (s *StorageDB) AddRefund(userID, orderID uint64, order *domain.Order) error {
	return s.txManager.RunReadCommitted(s.ctx, func(ctxTx context.Context) error {
		err := s.db.AddRefund(ctxTx, userID, orderID, order)
		if err != nil {
			return err
		}
		return s.db.SetOrderStatus(ctxTx, orderID, domain.StatusReturned)
	})
}

func (s *StorageDB) RemoveRefund(orderID uint64) error {
	return s.txManager.RunSerializable(s.ctx, func(ctxTx context.Context) error {
		err := s.db.RemoveRefund(ctxTx, orderID)
		if err != nil {
			return err
		}
		return s.db.SetOrderStatus(ctxTx, orderID, domain.StatusGiveCourier)
	})
}

func (s *StorageDB) GetRefunds(pageID, ordersPerPage uint64) (orders []domain.OrderView, err error) {
	err = s.txManager.RunReadOnlyCommitted(s.ctx, func(ctxTx context.Context) error {
		orders, err = s.db.GetRefunds(ctxTx, pageID, ordersPerPage)
		return err
	})
	return orders, err
}
