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

func newReturnUsecase(mocks *mocks) *ReturnUsecase {
	st := &storage_json.Storage{
		Ohp:   mocks.ohp,
		Rp:    mocks.rp,
		Users: mocks.up,
	}
	return NewReturnUsecase(st)
}

func TestReturnUsecase(t *testing.T) {
	type (
		args struct {
			req *dto.ReturnRequest
		}

		TestData struct {
			req         *dto.ReturnRequest
			orderStatus *domain.OrderStatus
		}
	)

	ctrl := minimock.NewController(t)
	m := newMocks(ctrl)
	u := newReturnUsecase(m)

	td := map[string]TestData{
		"SuccessReturned": {
			req: &dto.ReturnRequest{
				OrderID: 1,
			},
			orderStatus: &domain.OrderStatus{
				Status: domain.StatusReturned,
			},
		},
		"SuccessAccepted": {
			req: &dto.ReturnRequest{
				OrderID: 2,
			},
			orderStatus: &domain.OrderStatus{
				Status: domain.StatusAccepted,
				UserID: 2,
				Order: &domain.Order{
					ExpirationDate: "01-09-2024",
				},
			},
		},
		"ExpDateHasnotExpired": {
			req: &dto.ReturnRequest{
				OrderID: 3,
			},
			orderStatus: &domain.OrderStatus{
				Status: domain.StatusAccepted,
				UserID: 3,
				Order: &domain.Order{
					ExpirationDate: utils.CurrentDateString(),
				},
			},
		},
		"WrongOrderStatus": {
			req: &dto.ReturnRequest{
				OrderID: 4,
			},
			orderStatus: &domain.OrderStatus{
				Status: domain.StatusGiveClient,
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
			name: "SuccessReturned",
			args: args{td["SuccessReturned"].req},
			prepare: func() {
				data := td["SuccessReturned"]
				req := data.req
				stat := data.orderStatus

				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(stat, nil)
				m.rp.RemoveRefundMock.When(req.OrderID).Then(nil)
				m.ohp.SetOrderStatusMock.When(req.OrderID, domain.StatusGiveCourier).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "SuccessAccepted",
			args: args{td["SuccessAccepted"].req},
			prepare: func() {
				data := td["SuccessAccepted"]
				req := data.req
				stat := data.orderStatus

				expData, _ := time.Parse("02-01-2006", stat.ExpirationDate)
				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(stat, nil)
				m.up.GetExpirationDateMock.When(stat.UserID, req.OrderID).Then(expData, nil)
				m.up.RemoveOrderMock.When(stat.UserID, req.OrderID).Then(nil)
				m.ohp.SetOrderStatusMock.When(req.OrderID, domain.StatusGiveCourier).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "ExpDateHasnotExpired",
			args: args{td["ExpDateHasnotExpired"].req},
			prepare: func() {
				data := td["ExpDateHasnotExpired"]
				req := data.req
				stat := data.orderStatus

				expData, _ := time.Parse("02-01-2006", stat.ExpirationDate)
				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(stat, nil)
				m.up.GetExpirationDateMock.When(stat.UserID, req.OrderID).Then(expData, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "WrongOrderStatus",
			args: args{td["WrongOrderStatus"].req},
			prepare: func() {
				data := td["WrongOrderStatus"]
				req := data.req
				stat := data.orderStatus

				m.ohp.GetOrderStatusMock.When(req.OrderID).Then(stat, nil)
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

			err := u.Return(tt.args.req)
			tt.wantErr(t, err)
		})
	}
}
