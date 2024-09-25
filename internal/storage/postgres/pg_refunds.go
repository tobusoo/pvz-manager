package postgres

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

func (pg *PgRepository) AddRefund(ctx context.Context, userID, orderID uint64, order *domain.Order) error {
	return nil
}

func (pg *PgRepository) RemoveRefund(ctx context.Context, orderID uint64) error {
	return nil
}

func (pg *PgRepository) GetRefunds(ctx context.Context, pageID, ordersPerPage uint64) ([]domain.OrderView, error) {
	return nil, nil
}
