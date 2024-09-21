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
	st *storage.Storage
}

func NewGiveUsecase(st *storage.Storage) *GiveUsecase {
	return &GiveUsecase{st}
}

func (u *GiveUsecase) giveCheckErr(userID, orderID uint64, status *domain.OrderStatus) error {
	if status.UserID != userID {
		return fmt.Errorf("can't give order %d: different userID", orderID)
	}

	if status.Status != domain.StatusAccepted {
		return fmt.Errorf("can't give order %d: status = %s", orderID, status.Status)
	}

	expDate, err := u.st.GetExpirationDate(status.UserID, uint64(orderID))
	if err != nil {
		return fmt.Errorf("can't give: %s", err)
	}

	if utils.CurrentDate().After(expDate) {
		return fmt.Errorf("can't give order %d: expiration date has already passed", orderID)
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

func (u *GiveUsecase) giveProcess(orders []uint, isGoodResponse bool, errors []error) []error {
	if isGoodResponse {
		for _, order := range orders {
			u.st.RemoveOrder(uint64(order), domain.StatusGiveClient)
		}
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

	for _, orderID := range req.Orders {
		userID, knowUserID, err = u.giveCheckOrder(uint64(orderID), userID, knowUserID)
		if err != nil {
			errors = append(errors, err)
			isGoodResponse = false
		}
	}

	return u.giveProcess(req.Orders, isGoodResponse, errors)
}
