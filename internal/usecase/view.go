package usecase

import (
	"fmt"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

type ViewUsecase struct {
	st storage.Storage
}

func NewViewUsecase(st storage.Storage) *ViewUsecase {
	return &ViewUsecase{st}
}

func (u *ViewUsecase) GetRefunds(req *dto.ViewRefundsRequest) ([]domain.OrderView, error) {
	refunds, err := u.st.GetRefunds(req.PageID, req.OrdersPerPage)
	if err != nil {
		return nil, fmt.Errorf("error while view refund: %s", err)
	}

	if len(refunds) == 0 {
		return nil, fmt.Errorf("there are no refunds for page %d with ordersPerPage equal to %d: %w", req.PageID, req.OrdersPerPage, domain.ErrNotFound)
	}

	return refunds, nil
}

func (u *ViewUsecase) GetOrders(req *dto.ViewOrdersRequest) ([]domain.OrderView, error) {
	orders, err := u.st.GetOrdersByUserID(req.UserID, req.FirstOrderID, req.OrdersLimit)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("user %d doesn't have orders: %w", req.UserID, domain.ErrNotFound)
	}

	return orders, nil
}
