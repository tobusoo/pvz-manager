package usecase

import (
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/mock"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

type mocks struct {
	ohp *mock.OrdersHistoryRepositoryMock
	rp  *mock.RefundsRepositoryMock
	up  *mock.UsersRepositoryMock
}

func newMocks(ctrl *minimock.Controller) *mocks {
	return &mocks{
		ohp: mock.NewOrdersHistoryRepositoryMock(ctrl),
		rp:  mock.NewRefundsRepositoryMock(ctrl),
		up:  mock.NewUsersRepositoryMock(ctrl),
	}
}

func newAcceptUsecase(mocks *mocks) *AcceptUsecase {
	st := &storage.Storage{
		OrdersHistoryRepository: mocks.ohp,
		RefundsRepository:       mocks.rp,
		Users:                   mocks.up,
	}
	return NewAcceptUsecase(st)
}

func TestAcceptUsecase_AcceptOrder(t *testing.T) {
	type (
		args struct {
			req *dto.AddOrderRequest
		}

		TestData struct {
			req   *dto.AddOrderRequest
			order *domain.Order
		}
	)

	ctrl := minimock.NewController(t)

	td := map[string]TestData{
		"SuccessAccept": {
			req: &dto.AddOrderRequest{
				ExpirationDate: utils.CurrentDateString(),
				ContainerType:  "",
				UseTape:        false,
				UserID:         1,
				OrderID:        1,
				Cost:           100,
				Weight:         100,
			},
			order: &domain.Order{
				ExpirationDate: utils.CurrentDateString(),
				PackageType:    "default",
				Cost:           100,
				Weight:         100,
				UseTape:        false,
			},
		},
		"ExpirationDateWrongFormat": {
			req: &dto.AddOrderRequest{
				ExpirationDate: "10.12.2025",
			},
			order: nil,
		},
		"ExpiraionDatePassed": {
			req: &dto.AddOrderRequest{
				ExpirationDate: "10-10-2023",
			},
			order: nil,
		},
		"WrongContainerType": {
			req: &dto.AddOrderRequest{
				ExpirationDate: utils.CurrentDateString(),
				ContainerType:  "wrong",
			},
			order: nil,
		},
		"UsingTwoTape": {
			req: &dto.AddOrderRequest{
				ExpirationDate: utils.CurrentDateString(),
				ContainerType:  "tape",
				UseTape:        true,
			},
			order: nil,
		},
		"UsingTapeWithoutContainer": {
			req: &dto.AddOrderRequest{
				ExpirationDate: utils.CurrentDateString(),
				ContainerType:  "",
				UseTape:        true,
			},
			order: nil,
		},
	}

	tests := []struct {
		name    string
		args    args
		prepare func(m *mocks)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "SuccessAccept",
			args: args{td["SuccessAccept"].req},
			prepare: func(m *mocks) {
				data := td["SuccessAccept"]
				req := data.req
				order := data.order

				m.ohp.GetOrderStatusMock.Expect(req.UserID).Return(nil, fmt.Errorf("order %d not found", req.UserID))
				m.up.AddOrderMock.Expect(req.UserID, req.OrderID, order).Return(nil)
				m.ohp.AddOrderStatusMock.Expect(req.OrderID, req.UserID, domain.StatusAccepted, order).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "ExpirationDateWrongFormat",
			args: args{td["ExpirationDateWrongFormat"].req},
			prepare: func(m *mocks) {
			},
			wantErr: assert.Error,
		},
		{
			name: "ExpiraionDatePassed",
			args: args{td["ExpiraionDatePassed"].req},
			prepare: func(m *mocks) {
			},
			wantErr: assert.Error,
		},
		{
			name: "WrongContainerType",
			args: args{td["WrongContainerType"].req},
			prepare: func(m *mocks) {
			},
			wantErr: assert.Error,
		},
		{
			name: "UsingTwoTape",
			args: args{td["UsingTwoTape"].req},
			prepare: func(m *mocks) {
			},
			wantErr: assert.Error,
		},
		{
			name: "UsingTapeWithoutContainer",
			args: args{td["UsingTapeWithoutContainer"].req},
			prepare: func(m *mocks) {
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(ctrl)
			tt.prepare(mocks)
			u := newAcceptUsecase(mocks)

			err := u.AcceptOrder(tt.args.req)
			tt.wantErr(t, err)
		})
	}
}

func TestAcceptUsecase_AcceptRefund(t *testing.T) {
	type (
		args struct {
			req *dto.RefundRequest
		}

		TestData struct {
			req   *dto.RefundRequest
			order *domain.OrderStatus
		}
	)

	ctrl := minimock.NewController(t)

	td := map[string]TestData{
		"SuccessRefund": {
			req: &dto.RefundRequest{
				UserID:  1,
				OrderID: 1,
			},
			order: &domain.OrderStatus{
				Status: domain.StatusGiveClient,
				UserID: 1,
				Order: &domain.Order{
					ExpirationDate: utils.CurrentDateString(),
				},
				Date: utils.CurrentDateString(),
			},
		},
		"WrongOrderStatus": {
			req: &dto.RefundRequest{
				UserID:  2,
				OrderID: 2,
			},
			order: &domain.OrderStatus{
				Status: domain.StatusAccepted,
				UserID: 2,
			},
		},
		"WrongUserID": {
			req: &dto.RefundRequest{
				UserID:  3,
				OrderID: 3,
			},
			order: &domain.OrderStatus{
				Status: domain.StatusGiveClient,
				UserID: 4,
			},
		},
		"2DaysHavePassedSinceIssuedToClient": {
			req: &dto.RefundRequest{
				UserID:  5,
				OrderID: 5,
			},
			order: &domain.OrderStatus{
				Status: domain.StatusGiveClient,
				UserID: 5,
				Date:   "01-09-2024",
			},
		},
	}

	tests := []struct {
		name    string
		args    args
		prepare func(m *mocks)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "SuccessRefund",
			args: args{td["SuccessRefund"].req},
			prepare: func(m *mocks) {
				data := td["SuccessRefund"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.Expect(req.OrderID).Return(orderStat, nil)
				m.rp.AddRefundMock.Expect(req.UserID, req.OrderID, orderStat.Order).Return(nil)
				m.ohp.SetOrderStatusMock.Expect(req.OrderID, domain.StatusReturned).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "WrongOrderStatus",
			args: args{td["WrongOrderStatus"].req},
			prepare: func(m *mocks) {
				data := td["WrongOrderStatus"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.Expect(req.OrderID).Return(orderStat, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "WrongUserID",
			args: args{td["WrongUserID"].req},
			prepare: func(m *mocks) {
				data := td["WrongUserID"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.Expect(req.OrderID).Return(orderStat, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "2DaysHavePassedSinceIssuedToClient",
			args: args{td["2DaysHavePassedSinceIssuedToClient"].req},
			prepare: func(m *mocks) {
				data := td["2DaysHavePassedSinceIssuedToClient"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.Expect(req.OrderID).Return(orderStat, nil)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(ctrl)
			tt.prepare(mocks)
			u := newAcceptUsecase(mocks)

			err := u.AcceptRefund(tt.args.req)
			tt.wantErr(t, err)
		})
	}
}
