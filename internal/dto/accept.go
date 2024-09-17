package dto

type AddOrderRequest struct {
	ExpirationDate string `json:"expirationDate"`
	ContainerType  string `json:"containerType"`
	UserID         uint64 `json:"userID"`
	OrderID        uint64 `json:"orderID"`
	Cost           uint64 `json:"cost"`
	Weight         uint64 `json:"weight"`
	UseTape        bool   `json:"useTape"`
}

type RefundRequest struct {
	UserID  uint64 `json:"userID"`
	OrderID uint64 `json:"orderID"`
}
