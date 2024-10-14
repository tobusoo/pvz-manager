package usecase

import (
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/storage_json"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

func newGiveUsecase(mocks *mocks) *GiveUsecase {
	st := &storage_json.Storage{
		Ohp:   mocks.ohp,
		Rp:    mocks.rp,
		Users: mocks.up,
	}
	return NewGiveUsecase(st)
}

func TestGiveUsecase(t *testing.T) {
	type (
		args struct {
			req *dto.GiveOrdersRequest
		}

		TestData struct {
			req         *dto.GiveOrdersRequest
			orderStatus []*domain.OrderStatus
		}
	)

	ctrl := minimock.NewController(t)
	m := newMocks(ctrl)
	u := newGiveUsecase(m)

	td := map[string]TestData{
		"SuccessGive": {
			req: &dto.GiveOrdersRequest{
				Orders: []uint64{1, 2},
			},
			orderStatus: []*domain.OrderStatus{
				{
					UserID: 1,
					Status: domain.StatusAccepted,
					Order: &domain.Order{
						ExpirationDate: utils.CurrentDateString(),
					},
				},
				{
					UserID: 1,
					Status: domain.StatusAccepted,
					Order: &domain.Order{
						ExpirationDate: utils.CurrentDateString(),
					},
				},
			},
		},
		"DifferentUserID": {
			req: &dto.GiveOrdersRequest{
				Orders: []uint64{3, 4},
			},
			orderStatus: []*domain.OrderStatus{
				{
					UserID: 2,
					Status: domain.StatusAccepted,
					Order: &domain.Order{
						ExpirationDate: utils.CurrentDateString(),
					},
				},
				{
					UserID: 3,
					Status: domain.StatusAccepted,
					Order: &domain.Order{
						ExpirationDate: utils.CurrentDateString(),
					},
				},
			},
		},
		"BadOrderStatus": {
			req: &dto.GiveOrdersRequest{
				Orders: []uint64{5},
			},
			orderStatus: []*domain.OrderStatus{
				{
					UserID: 5,
					Status: domain.StatusGiveClient,
					Order:  nil,
				},
			},
		},
		"ExpirationDatePassed": {
			req: &dto.GiveOrdersRequest{
				Orders: []uint64{6},
			},
			orderStatus: []*domain.OrderStatus{
				{
					UserID: 6,
					Status: domain.StatusAccepted,
					Order: &domain.Order{
						ExpirationDate: "18-09-2024",
					},
				},
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
			name: "SuccessGive",
			args: args{td["SuccessGive"].req},
			prepare: func() {
				data := td["SuccessGive"]
				ordersStat := data.orderStatus
				orders := data.req.Orders

				for i, order := range orders {
					stat := ordersStat[i]
					orderID := uint64(order)

					m.ohp.GetOrderStatusMock.When(orderID).Then(stat, nil)
					m.up.GetExpirationDateMock.When(stat.UserID, orderID).Then(utils.CurrentDate(), nil)
					m.up.CanRemoveMock.When(stat.UserID, orderID).Then(nil)
					m.up.RemoveOrderMock.When(stat.UserID, orderID).Then(nil)
					m.ohp.SetOrderStatusMock.When(orderID, domain.StatusGiveClient).Then(nil)
				}
			},
			wantErr: assert.NoError,
		},
		{
			name: "DifferentUserID",
			args: args{td["DifferentUserID"].req},
			prepare: func() {
				data := td["DifferentUserID"]
				ordersStat := data.orderStatus
				orders := data.req.Orders

				for i, order := range orders {
					stat := ordersStat[i]
					orderID := uint64(order)

					m.ohp.GetOrderStatusMock.Optional().When(orderID).Then(stat, nil)
					m.up.GetExpirationDateMock.Optional().When(stat.UserID, orderID).Then(utils.CurrentDate(), nil)
					m.up.CanRemoveMock.Optional().When(stat.UserID, orderID).Then(nil)
					m.up.RemoveOrderMock.Optional().When(stat.UserID, orderID).Then(nil)
					m.ohp.SetOrderStatusMock.Optional().When(orderID, domain.StatusGiveClient).Then(nil)
				}
			},
			wantErr: assert.Error,
		},
		{
			name: "BadOrderStatus",
			args: args{td["BadOrderStatus"].req},
			prepare: func() {
				data := td["BadOrderStatus"]
				stat := data.orderStatus[0]
				orderID := uint64(data.req.Orders[0])

				m.ohp.GetOrderStatusMock.When(orderID).Then(stat, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "ExpirationDatePassed",
			args: args{td["ExpirationDatePassed"].req},
			prepare: func() {
				data := td["ExpirationDatePassed"]
				stat := data.orderStatus[0]
				orderID := uint64(data.req.Orders[0])

				expData, _ := time.Parse("02-01-2006", stat.ExpirationDate)
				m.ohp.GetOrderStatusMock.When(orderID).Then(stat, nil)
				m.up.GetExpirationDateMock.When(stat.UserID, orderID).Then(expData, nil)
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

			var err error
			errors := u.Give(tt.args.req)
			if len(errors) > 0 {
				err = errors[0]
			}

			tt.wantErr(t, err)
		})
	}
}
