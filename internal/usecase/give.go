package usecase

import (
	"fmt"
	"slices"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

type GiveUsecase struct {
	st storage.Storage
}

func NewGiveUsecase(st storage.Storage) *GiveUsecase {
	return &GiveUsecase{st}
}

func (u *GiveUsecase) giveCheckErr(userID, orderID uint64, status *domain.OrderStatus) error {
	if status.UserID != userID {
		return fmt.Errorf("can't give order %d: different userID: %w", orderID, domain.ErrWrongInput)
	}

	if status.Status != domain.StatusAccepted {
		return fmt.Errorf("can't give order %d: status = %s: %w", orderID, status.Status, domain.ErrWrongStatus)
	}

	expDate, err := u.st.GetExpirationDate(status.UserID, uint64(orderID))
	if err != nil {
		return fmt.Errorf("can't give: %s", err)
	}

	if utils.CurrentDate().After(expDate) {
		return fmt.Errorf("can't give order %d: %w", orderID, domain.ErrExpirationDatePassed)
	}

	return u.st.CanRemoveOrder(orderID)
}

func (u *GiveUsecase) giveCheckOrder(orderID, userID uint64, knowUserID bool) (uint64, bool, error) {
	status, err := u.st.GetOrderStatus(orderID)
	if err != nil {
		return userID, knowUserID, fmt.Errorf("can't give: %s", err)
	}

	if !knowUserID {
		userID = status.UserID
		knowUserID = true
	}

	return userID, knowUserID, u.giveCheckErr(userID, orderID, status)
}

func (u *GiveUsecase) giveProcess(orders []uint64, isGoodResponse bool, errors []error) []error {
	if isGoodResponse {
		u.st.RemoveOrders(orders, domain.StatusGiveClient)
		return nil
	}

	return errors
}

func (u *GiveUsecase) Give(req *dto.GiveOrdersRequest) []error {
	var err error
	userID := uint64(0)
	knowUserID := false
	isGoodResponse := true
	errors := make([]error, 0)

	slices.Sort(req.Orders)
	req.Orders = slices.Compact(req.Orders)
	orders := make([]uint64, 0, len(req.Orders))

	for _, orderID := range req.Orders {
		userID, knowUserID, err = u.giveCheckOrder(uint64(orderID), userID, knowUserID)
		if err != nil {
			errors = append(errors, err)
			isGoodResponse = false
		}
		orders = append(orders, uint64(orderID))
	}

	return u.giveProcess(orders, isGoodResponse, errors)
}
