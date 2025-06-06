package usecase

import (
	"fmt"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

type AcceptUsecase struct {
	st storage.Storage
}

func NewAcceptUsecase(st storage.Storage) *AcceptUsecase {
	return &AcceptUsecase{st}
}

func addAdditionalTape(req *dto.AddOrderRequest, cs strategy.ContainerStrategy) error {
	if req.ContainerType == "" {
		return fmt.Errorf("can't use additional tape: containerType must be defined: %w", domain.ErrWrongInput)
	}

	return cs.UseTape()
}

func generateOrder(req *dto.AddOrderRequest) (*domain.Order, error) {
	var cs strategy.ContainerStrategy
	cs, ok := strategy.ContainerTypeMap[req.ContainerType]
	if !ok {
		return nil, fmt.Errorf("%s isn't container type: %w", req.ContainerType, domain.ErrWrongInput)
	}

	if req.UseTape {
		err := addAdditionalTape(req, cs)
		if err != nil {
			return nil, err
		}
	}

	return domain.NewOrder(req.Cost, req.Weight, req.ExpirationDate, cs)
}

func (u *AcceptUsecase) AcceptOrder(req *dto.AddOrderRequest) error {
	expDate, err := time.Parse("02-01-2006", req.ExpirationDate)
	if err != nil {
		return fmt.Errorf("%w: %w", err, domain.ErrWrongInput)
	}

	currentDate := utils.CurrentDate()
	if currentDate.After(expDate) {
		return domain.ErrExpirationDatePassed
	}

	req.ExpirationDate = expDate.Format("02-01-2006")
	order, err := generateOrder(req)
	if err != nil {
		return err
	}

	return u.st.AddOrder(req.UserID, req.OrderID, order)
}

func acceptRefundCheckErr(req *dto.RefundRequest, order *domain.OrderStatus) error {
	if order.Status != domain.StatusGiveClient {
		return fmt.Errorf("can not refund order %d: status = %s: %w", req.OrderID, order.Status, domain.ErrWrongStatus)
	}

	if req.UserID != order.UserID {
		return fmt.Errorf("can not refund order %d: wrong userID: %w", req.OrderID, domain.ErrWrongInput)
	}

	issuedDate, err := time.Parse("02-01-2006", order.UpdatedAt)
	if err != nil {
		return err
	}

	issuedDate = issuedDate.Add(2 * 24 * time.Hour)
	currentDate := utils.CurrentDate()

	if currentDate.After(issuedDate) {
		return domain.ErrTwoDaysPassed
	}

	return nil
}

func (u *AcceptUsecase) AcceptRefund(req *dto.RefundRequest) error {
	order, err := u.st.GetOrderStatus(req.OrderID)
	if err != nil {
		return err
	}

	if err = acceptRefundCheckErr(req, order); err != nil {
		return err
	}

	return u.st.AddRefund(req.UserID, req.OrderID, order.Order)
}
