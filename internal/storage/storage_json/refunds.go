package storage_json

import (
	"fmt"
	"sync"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
)

type Refunds struct {
	Orders          []domain.OrderView `json:"orders"`
	OrdersIDatArray map[uint64]int     `json:"ordersIDatArray"`

	mtx sync.Mutex
}

func NewRefunds() *Refunds {
	return &Refunds{
		Orders:          make([]domain.OrderView, 0),
		OrdersIDatArray: make(map[uint64]int),
	}
}

func (r *Refunds) AddRefund(userID, orderID uint64, order *domain.Order) (err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.Orders = append(r.Orders, domain.OrderView{
		Order:   order,
		UserID:  userID,
		OrderID: orderID,
		Exist:   true,
	})
	r.OrdersIDatArray[orderID] = len(r.Orders) - 1

	return nil
}

func (r *Refunds) RemoveRefund(orderID uint64) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	id, ok := r.OrdersIDatArray[orderID]
	if !ok {
		return fmt.Errorf("not found order %d at refunded", orderID)
	}

	r.Orders[id].Exist = false
	return nil
}

func (r *Refunds) GetRefunds(pageID, ordersPerPage uint64) (res []domain.OrderView, err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if err := r.getRefundsCheckErr(pageID, ordersPerPage); err != nil {
		return nil, err
	}

	firstOrderID := int((pageID - 1) * ordersPerPage)
	return r.getRefundsSlice(firstOrderID, ordersPerPage), nil
}

func (r *Refunds) getRefundsSlice(firstOrderID int, ordersPerPage uint64) []domain.OrderView {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	res := make([]domain.OrderView, 0)
	ordersCount := 0

	for i := firstOrderID; i < len(r.Orders); i++ {
		if r.Orders[i].Exist {
			res = append(res, r.Orders[i])
			ordersCount++
		}

		if ordersCount == int(ordersPerPage) {
			break
		}
	}

	return res
}

func (r *Refunds) getRefundsCheckErr(pageID, ordersPerPage uint64) error {
	if ordersPerPage == 0 {
		return fmt.Errorf("orders per page must be greater than 0")
	}

	if pageID == 0 {
		return fmt.Errorf("pageID myst be greater than 0")
	}

	return nil
}
