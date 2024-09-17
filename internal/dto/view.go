package dto

type ViewRefundsRequest struct {
	PageID        uint64 `json:"pageID"`
	OrdersPerPage uint64 `json:"ordersPerPage"`
}

type ViewOrdersRequest struct {
	UserID       uint64 `json:"userID"`
	FirstOrderID uint64 `json:"firstOrderID"`
	OrdersLimit  uint64 `json:"ordersLimit"`
}
