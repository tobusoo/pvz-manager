package dto

type GiveOrdersRequest struct {
	Orders []uint64 `json:"orders"`
}
