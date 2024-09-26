package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

func (pg *PgRepository) AddRefund(ctx context.Context, userID, orderID uint64, order *domain.Order) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		insert into refunds(
		order_id)
		values $1`,
		orderID,
	)

	return err
}

func (pg *PgRepository) RemoveRefund(ctx context.Context, orderID uint64) error {
	tx := pg.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		delete from refunds
		where order_id = $1`,
		orderID,
	)

	return err
}

func (pg *PgRepository) GetRefunds(ctx context.Context, pageID, ordersPerPage uint64) ([]domain.OrderView, error) {
	var orders []domain.OrderView

	tx := pg.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, `
		select
			oh.user_id,
			oh.order_id,
			oh.expiration_date,
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
		) r on oh.order_id = r.order_id`,
		ordersPerPage,
		(pageID-1)*ordersPerPage,
	)

	return orders, err
}
