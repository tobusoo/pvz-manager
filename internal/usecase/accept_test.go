package usecase

import (
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/storage_json"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/storage_json/mock"
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
	st := &storage_json.Storage{
		Ohp:   mocks.ohp,
		Rp:    mocks.rp,
		Users: mocks.up,
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

	m := newMocks(ctrl)
	u := newAcceptUsecase(m)

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
		prepare func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "SuccessAccept",
			args: args{td["SuccessAccept"].req},
			prepare: func() {
				data := td["SuccessAccept"]
				req := data.req
				order := data.order

				m.ohp.GetOrderStatusMock.When(req.UserID).Then(nil, fmt.Errorf("order %d not found", req.UserID))
				m.up.AddOrderMock.When(req.UserID, req.OrderID, order).Then(nil)
				m.ohp.AddOrderStatusMock.When(req.OrderID, req.UserID, domain.StatusAccepted, order).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "ExpirationDateWrongFormat",
			args: args{td["ExpirationDateWrongFormat"].req},
			prepare: func() {
			},
			wantErr: assert.Error,
		},
		{
			name: "ExpiraionDatePassed",
			args: args{td["ExpiraionDatePassed"].req},
			prepare: func() {
			},
			wantErr: assert.Error,
		},
		{
			name: "WrongContainerType",
			args: args{td["WrongContainerType"].req},
			prepare: func() {
			},
			wantErr: assert.Error,
		},
		{
			name: "UsingTwoTape",
			args: args{td["UsingTwoTape"].req},
			prepare: func() {
			},
			wantErr: assert.Error,
		},
		{
			name: "UsingTapeWithoutContainer",
			args: args{td["UsingTapeWithoutContainer"].req},
			prepare: func() {
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt.prepare()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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
	m := newMocks(ctrl)
	u := newAcceptUsecase(m)

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
				UpdatedAt: utils.CurrentDateString(),
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
				Status:    domain.StatusGiveClient,
				UserID:    5,
				UpdatedAt: "01-09-2024",
			},
		},
	}

	tests := []struct {
		name    string
		args    args
		prepare func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "SuccessRefund",
			args: args{td["SuccessRefund"].req},
			prepare: func() {
				data := td["SuccessRefund"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(orderStat, nil)
				m.rp.AddRefundMock.When(req.UserID, req.OrderID, orderStat.Order).Then(nil)
				m.ohp.SetOrderStatusMock.When(req.OrderID, domain.StatusReturned).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "WrongOrderStatus",
			args: args{td["WrongOrderStatus"].req},
			prepare: func() {
				data := td["WrongOrderStatus"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(orderStat, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "WrongUserID",
			args: args{td["WrongUserID"].req},
			prepare: func() {
				data := td["WrongUserID"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(orderStat, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "2DaysHavePassedSinceIssuedToClient",
			args: args{td["2DaysHavePassedSinceIssuedToClient"].req},
			prepare: func() {
				data := td["2DaysHavePassedSinceIssuedToClient"]
				req := data.req
				orderStat := data.order

				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(orderStat, nil)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt.prepare()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := u.AcceptRefund(tt.args.req)
			tt.wantErr(t, err)
		})
	}
}
