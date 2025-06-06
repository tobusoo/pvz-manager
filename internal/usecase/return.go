package usecase

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

type ReturnUsecase struct {
	st storage.Storage
}

func NewReturnUsecase(st storage.Storage) *ReturnUsecase {
	return &ReturnUsecase{st}
}

func (u *ReturnUsecase) returnAccepted(orderID uint64, order *domain.OrderStatus) error {
	expDate, err := u.st.GetExpirationDate(order.UserID, orderID)
	if err != nil {
		return err
	}

	if expDate.Add(24 * time.Hour).After(utils.CurrentDate()) {
		return fmt.Errorf("can't return order %d: %w", orderID, domain.ErrNotExpirationDate)
	}

	return u.st.RemoveOrder(orderID, domain.StatusGiveCourier)
}

func (u *ReturnUsecase) Return(req *dto.ReturnRequest) error {
	order, err := u.st.GetOrderStatus(req.OrderID)
	if err != nil {
		return err
	}

	switch order.Status {
	case domain.StatusReturned:
		u.st.RemoveRefund(req.OrderID)
		u.st.SetOrderStatus(req.OrderID, domain.StatusGiveCourier)

	case domain.StatusAccepted:
		return u.returnAccepted(req.OrderID, order)
	default:
		return fmt.Errorf("can't return order %d: status = %s: %w", req.OrderID, order.Status, domain.ErrWrongStatus)
	}

	return nil
}
