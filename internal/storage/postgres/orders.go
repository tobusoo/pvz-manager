package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

func (pg *PgRepository) AddOrder(ctx context.Context, userID, orderID uint64) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `insert into orders(
		user_id,
		order_id)
		values ($1, $2)`,
		userID,
		orderID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("order %d already exists", orderID)
		}
		return fmt.Errorf("AddOrder: %w", err)
	}

	return nil
}

func (pg *PgRepository) GetOrder(ctx context.Context, userID, orderID uint64) (*domain.Order, error) {
	var order domain.Order
	tx := pg.txManager.GetQueryEngine(ctx)

	err := pgxscan.Get(ctx, tx, &order, `
		select 
			to_char(expiration_date, 'DD-MM-YYYY') as expiration_date,
			package_type,
			cost,
			weight,
			use_tape
		from orders_history
		where order_id = $1 and user_id = $2`,
		orderID,
		userID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("order %d for user %d not found", orderID, userID)
	} else if err != nil {
		return nil, fmt.Errorf("GetOrder: %w", err)
	}

	return &order, nil
}

func (pg *PgRepository) GetExpirationDate(ctx context.Context, userID, orderID uint64) (time.Time, error) {
	var expDate string
	tx := pg.txManager.GetQueryEngine(ctx)

	err := pgxscan.Get(ctx, tx, &expDate, `
		select 
			to_char(expiration_date, 'DD-MM-YYYY') as expiration_date
		from orders_history
		where user_id = $2 and order_id = $1`,
		orderID,
		userID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, fmt.Errorf("order %d for user %d not found", orderID, userID)
	} else if err != nil {
		return time.Time{}, fmt.Errorf("GetExpirationDate: %w", err)
	}

	return utils.StringToTime(expDate)
}

func (pg *PgRepository) GetOrdersByUserID(ctx context.Context, userID, firstOrderID, limit uint64) ([]domain.OrderView, error) {
	var orders []domain.OrderView
	tx := pg.txManager.GetQueryEngine(ctx)

	if err := pgxscan.Select(ctx, tx, &orders, `
		select
			user_id, 
			order_id,
			to_char(expiration_date, 'DD-MM-YYYY') as expiration_date,
			package_type,
			cost,
			weight,
			use_tape
		from orders_history
		where user_id = $1 and order_id >= $2 order by order_id limit $3`,
		userID,
		firstOrderID,
		limit,
	); err != nil {
		return nil, fmt.Errorf("GetOrdersByUserID: %w", err)
	}

	return orders, nil
}

func (pg *PgRepository) CanRemoveOrder(ctx context.Context, userID, orderID uint64) error {
	var exist bool
	tx := pg.txManager.GetQueryEngine(ctx)

	if err := pgxscan.Get(ctx, tx, &exist, `
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

	if !exist {
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
