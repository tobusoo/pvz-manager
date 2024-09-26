package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
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
	var order []*domain.Order
	tx := pg.txManager.GetQueryEngine(ctx)

	err := pgxscan.Select(ctx, tx, &order, `
		select 
			expiration_date,
			package_type,
			cost,
			weight,
			use_tape,
		from orders
		where user_id = $1 and order_id = $2`,
		userID,
		orderID,
	)

	if len(order) == 0 || err != nil {
		return nil, fmt.Errorf("not found user %d order's %d: %w", userID, orderID, err)
	}

	return order[0], nil
}

func (pg *PgRepository) GetExpirationDate(ctx context.Context, userID, orderID uint64) (time.Time, error) {
	var expDate []string
	tx := pg.txManager.GetQueryEngine(ctx)

	if err := pgxscan.Select(ctx, tx, &expDate, `
		select 
			expiration_date
		from orders
		where user_id = $1 and order_id = $2`,
		userID,
		orderID,
	); err != nil {
		return time.Time{}, fmt.Errorf("GetExpitarionDate: %w", err)
	}

	if len(expDate) == 0 {
		return time.Time{}, fmt.Errorf("not found expiration date for user %d order %d", userID, orderID)
	}

	return utils.StringToTime(expDate[0])
}

func (pg *PgRepository) GetOrdersByUserID(ctx context.Context, userID, firstOrderID, limit uint64) ([]domain.OrderView, error) {
	var orders []domain.OrderView
	tx := pg.txManager.GetQueryEngine(ctx)

	err := pgxscan.Select(ctx, tx, &orders, `
		select 
			order_id,
			expiration_date,
			package_type,
			cost,
			weight,
			use_tape
		from orders
		where user_id = $1 and order_id >= $2 order by order_id limit $3`,
		userID,
		firstOrderID,
		limit,
	)

	return orders, err
}

func (pg *PgRepository) CanRemoveOrder(ctx context.Context, userID, orderID uint64) error {
	var exist []bool
	tx := pg.txManager.GetQueryEngine(ctx)

	if err := pgxscan.Select(ctx, tx, &exist, `
		select exists (
			select 1
			from orders
			where user_id = $1 and order_id = $2
			)`,
		userID,
		orderID,
	); err != nil {
		return fmt.Errorf("CanRemoveOrder: %w", err)
	}

	if !exist[0] {
		return fmt.Errorf("not found user %d order's %d", userID, orderID)
	}

	return nil
}

func (pg *PgRepository) RemoveOrder(ctx context.Context, userID, orderID uint64) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	result, err := tx.Exec(ctx, `
		delete from orders
		where user_id = $1 and order_id = $2
		`,
		userID,
		orderID,
	)

	if err != nil {
		return fmt.Errorf("RemoveOrder: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user %d order's %d not found", userID, orderID)
	}

	return err
}
