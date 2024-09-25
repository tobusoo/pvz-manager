package postgres

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

func (pg *PgRepository) AddOrderStatus(ctx context.Context, orderID, userID uint64, status string, order *domain.Order) error {
	return nil
}

func (pg *PgRepository) GetOrderStatus(ctx context.Context, orderID uint64) (*domain.OrderStatus, error) {
	return nil, fmt.Errorf("just for testing")
}

func (pg *PgRepository) SetOrderStatus(ctx context.Context, orderID uint64, status string) error {
	return nil
}
