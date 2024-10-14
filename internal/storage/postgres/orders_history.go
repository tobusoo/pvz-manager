package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

func (pg *PgRepository) AddOrderStatus(ctx context.Context, orderID, userID uint64, status string, order *domain.Order) error {
	tx := pg.txManager.GetQueryEngine(ctx)
	expDate, err := utils.StringToTime(order.ExpirationDate)
	if err != nil {
		return fmt.Errorf("AddOrderStatus: %w", err)
	}

	_, err = tx.Exec(ctx,
		`insert into orders_history(
		order_id,
		user_id,
		expiration_date,
		package_type,
		weight,
		cost,
		use_tape,
		status,
		updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		orderID,
		userID,
		expDate,
		order.PackageType,
		order.Weight,
		order.Cost,
		order.UseTape,
		status,
		utils.CurrentDate(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExist
		}
		return fmt.Errorf("AddOrderStatus: %w", err)
	}

	return err
}

func (pg *PgRepository) GetOrderOnlyStatus(ctx context.Context, orderID uint64) (string, error) {
	var status string

	tx := pg.txManager.GetQueryEngine(ctx)
	err := pgxscan.Get(ctx, tx, &status,
		`select 
		 	status
		 from orders_history
		 where order_id = $1`,
		orderID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("GetOrderOnlyStatus: %w", err)
	}

	return status, nil
}

func (pg *PgRepository) GetOrderStatus(ctx context.Context, orderID uint64) (*domain.OrderStatus, error) {
	var order domain.OrderStatus

	tx := pg.txManager.GetQueryEngine(ctx)
	err := pgxscan.Get(ctx, tx, &order,
		`select 
		 user_id,
		 to_char(expiration_date, 'DD-MM-YYYY') as expiration_date,
		 package_type,
		 weight,
		 cost,
		 use_tape,
		 status,
		 to_char(updated_at, 'DD-MM-YYYY') as updated_at
		 from orders_history
		 where order_id = $1`,
		orderID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("GetOrderStatus: %w", err)
	}

	return &order, nil
}

func (pg *PgRepository) SetOrderStatus(ctx context.Context, orderID uint64, status string) error {
	tx := pg.txManager.GetQueryEngine(ctx)
	result, err := tx.Exec(ctx,
		`update orders_history
		 set status = $2
      	 where order_id = $1`,
		orderID,
		status,
	)

	if err != nil {
		return fmt.Errorf("SetOrderStatus: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
