package postgres

import (
	"context"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

func (pg *PgRepository) AddOrder(ctx context.Context, userID, orderID uint64, order *domain.Order) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `insert into orders(
		user_id,
		order_id,
		expiration_date,
		package_type,
		weight,
		cost,
		use_tape)
		values ($1, $2, $3, $4, $5, $6, $7)
	`, userID,
		orderID,
		order.ExpirationDate,
		order.PackageType,
		order.Weight,
		order.Cost,
		order.UseTape,
	)

	return err
}

func (pg *PgRepository) GetOrder(ctx context.Context, userID, orderID uint64) (*domain.Order, error) {
	return nil, nil
}

func (pg *PgRepository) GetExpirationDate(ctx context.Context, userID, orderID uint64) (time.Time, error) {
	return time.Time{}, nil
}

func (pg *PgRepository) GetOrdersByUserID(ctx context.Context, userID, firstOrderID, limit uint64) ([]domain.OrderView, error) {
	return nil, nil
}

func (pg *PgRepository) CanRemoveOrder(ctx context.Context, orderID uint64) error {
	return nil
}

func (pg *PgRepository) RemoveOrder(ctx context.Context, orderID uint64, status string) error {
	return nil
}

func (pg *PgRepository) RemoveOrders(ctx context.Context, ordersID []uint64, status string) error {
	return nil
}
