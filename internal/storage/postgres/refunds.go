package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

func (pg *PgRepository) AddRefund(ctx context.Context, userID, orderID uint64, order *domain.Order) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		insert into refunds(
			order_id)
		values ($1)`,
		orderID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExist
		}
		return fmt.Errorf("AddRefund: %w", err)
	}

	return nil
}

func (pg *PgRepository) RemoveRefund(ctx context.Context, orderID uint64) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	result, err := tx.Exec(ctx, `
		delete from refunds
		where order_id = $1`,
		orderID,
	)

	if err != nil {
		return fmt.Errorf("RemoveRefund: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("refund: %w", domain.ErrNotFound)
	}

	return err
}

func (pg *PgRepository) GetRefunds(ctx context.Context, pageID, ordersPerPage uint64) ([]domain.OrderView, error) {
	var orders []domain.OrderView
	limit := (pageID - 1) * ordersPerPage

	tx := pg.txManager.GetQueryEngine(ctx)
	if err := pgxscan.Select(ctx, tx, &orders, `
		select
			oh.user_id,
			oh.order_id,
			to_char(oh.expiration_date, 'DD-MM-YYYY') as expiration_date,
			oh.package_type,
			oh.weight,
			oh.cost,
			oh.use_tape
		from orders_history oh
		join (
			select order_id
			from refunds
			order by order_id
			limit $1 offset $2
		) r on oh.order_id = r.order_id
		order by oh.order_id`,
		ordersPerPage,
		limit,
	); err != nil {
		return nil, fmt.Errorf("GetRefunds: %w", err)
	}

	return orders, nil
}
