package clients

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
)

type ManagerService interface {
	AddOrder(ctx context.Context, req *dto.AddOrderRequest) error
	Refund(ctx context.Context, req *dto.RefundRequest) error
	GiveOrders(ctx context.Context, req *dto.GiveOrdersRequest) error
	Return(ctx context.Context, req *dto.ReturnRequest) error
	ViewOrders(ctx context.Context, req *dto.ViewOrdersRequest) (*dto.ViewOrdersResponse, error)
	ViewRefunds(ctx context.Context, req *dto.ViewRefundsRequest) (*dto.ViewRefundsResponse, error)
}
