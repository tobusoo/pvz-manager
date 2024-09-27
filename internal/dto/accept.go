package dto

type AddOrderRequest struct {
	ExpirationDate string `json:"expirationDate" fake:"{datefuture}"`
	ContainerType  string `json:"containerType" fake:"{randomstring:[, tape, box, package]}"`
	UserID         uint64 `json:"userID" fake:"{number:1,9223372036854775807}"`
	OrderID        uint64 `json:"orderID" fake:"{number:1,9223372036854775807}"`
	Cost           uint64 `json:"cost" fake:"{number:1,250000}"`
	Weight         uint64 `json:"weight" fake:"skip"`
	UseTape        bool   `json:"useTape"`
}

type RefundRequest struct {
	UserID  uint64 `json:"userID"`
	OrderID uint64 `json:"orderID"`
}
