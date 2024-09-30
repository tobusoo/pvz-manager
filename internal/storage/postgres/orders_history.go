package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

func (pg *PgRepository) AddOrderStatus(ctx context.Context, orderID, userID uint64, status string, order *domain.Order) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	if _, err := tx.Exec(ctx,
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
		order.ExpirationDate,
		order.PackageType,
		order.Weight,
		order.Cost,
		order.UseTape,
		status,
		utils.CurrentDateString(),
	); err != nil {
		return fmt.Errorf("AddOrderStatus: %w", err)
	}

	return nil
}

func (pg *PgRepository) GetOrderOnlyStatus(ctx context.Context, orderID uint64) (string, error) {
	var statuses []string

	tx := pg.txManager.GetQueryEngine(ctx)
	if err := pgxscan.Select(ctx, tx, &statuses,
		`select 
		 status,
		 from orders_history
		 where order_id = $1`,
		orderID,
	); err != nil {
		return "", fmt.Errorf("GetOrderOnlyStatus: %w", err)
	}

	if len(statuses) == 0 {
		return "", fmt.Errorf("not found order %d", orderID)
	}

	return statuses[0], nil

}

func (pg *PgRepository) GetOrderStatus(ctx context.Context, orderID uint64) (*domain.OrderStatus, error) {
	var orders []*domain.OrderStatus

	tx := pg.txManager.GetQueryEngine(ctx)
	if err := pgxscan.Select(ctx, tx, &orders,
		`select 
		 user_id,
		 expiration_date,
		 package_type,
		 weight,
		 cost,
		 use_tape,
		 status,
		 updated_at
		 from orders_history
		 where order_id = $1`,
		orderID,
	); err != nil {
		return nil, fmt.Errorf("GetOrderStatus: %w", err)
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("not found order %d", orderID)
	}

	return orders[0], nil
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
		return fmt.Errorf("order %d not found", orderID)
	}

	return nil
}
