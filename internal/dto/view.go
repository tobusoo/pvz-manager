package dto

import (
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type ViewRefundsRequest struct {
	PageID        uint64 `json:"pageID"`
	OrdersPerPage uint64 `json:"ordersPerPage"`
}

type ViewRefundsResponse struct {
	Orders []domain.OrderView
}

type ViewOrdersRequest struct {
	UserID       uint64 `json:"userID"`
	FirstOrderID uint64 `json:"firstOrderID"`
	OrdersLimit  uint64 `json:"ordersLimit"`
}

type ViewOrdersResponse struct {
	Orders []domain.OrderView
}
