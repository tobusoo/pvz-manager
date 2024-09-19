package usecase

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

func newViewUsecase(mocks *mocks) *ViewUsecase {
	st := &storage.Storage{
		OrdersHistoryRepository: mocks.ohp,
		RefundsRepository:       mocks.rp,
		Users:                   mocks.up,
	}
	return NewViewUsecase(st)
}

func TestViewUsecase_GetOrders(t *testing.T) {
	type (
		args struct {
			req    *dto.ViewOrdersRequest
			expect []domain.OrderView
		}

		TestData struct {
			req  *dto.ViewOrdersRequest
			view []domain.OrderView
		}
	)

	ctrl := minimock.NewController(t)
	m := newMocks(ctrl)
	u := newViewUsecase(m)

	td := map[string]TestData{
		"Success": {
			req: &dto.ViewOrdersRequest{
				UserID:       1,
				FirstOrderID: 1,
				OrdersLimit:  5,
			},
			view: []domain.OrderView{
				{
					OrderID: 1,
				},
				{
					OrderID: 2,
				},
				{
					OrderID: 3,
				},
			},
		},
		"ZeroLen": {
			req: &dto.ViewOrdersRequest{
				UserID:       2,
				FirstOrderID: 1,
				OrdersLimit:  5,
			},
			view: nil,
		},
	}

	tests := []struct {
		name    string
		args    args
		prepare func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{td["Success"].req, td["Success"].view},
			prepare: func() {
				data := td["Success"]
				req := data.req
				orders := data.view

				m.up.GetOrdersMock.When(req.UserID, req.FirstOrderID, req.OrdersLimit).Then(orders, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "ZeroLen",
			args: args{td["ZeroLen"].req, td["ZeroLen"].view},
			prepare: func() {
				data := td["ZeroLen"]
				req := data.req
				orders := data.view

				m.up.GetOrdersMock.When(req.UserID, req.FirstOrderID, req.OrdersLimit).Then(orders, nil)
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

			got, err := u.GetOrders(tt.args.req)
			tt.wantErr(t, err)

			require.Equal(t, tt.args.expect, got)
		})
	}
}

func TestViewUsecase_GetRefunds(t *testing.T) {
	type (
		args struct {
			req    *dto.ViewRefundsRequest
			expect []domain.OrderView
		}

		TestData struct {
			req  *dto.ViewRefundsRequest
			view []domain.OrderView
		}
	)

	ctrl := minimock.NewController(t)
	m := newMocks(ctrl)
	u := newViewUsecase(m)

	td := map[string]TestData{
		"Success": {
			req: &dto.ViewRefundsRequest{
				PageID:        1,
				OrdersPerPage: 10,
			},
			view: []domain.OrderView{
				{
					OrderID: 1,
				},
				{
					OrderID: 2,
				},
				{
					OrderID: 3,
				},
			},
		},
		"ZeroLen": {
			req: &dto.ViewRefundsRequest{
				PageID:        100,
				OrdersPerPage: 10,
			},
			view: nil,
		},
	}

	tests := []struct {
		name    string
		args    args
		prepare func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{td["Success"].req, td["Success"].view},
			prepare: func() {
				data := td["Success"]
				req := data.req
				orders := data.view

				m.rp.GetRefundsMock.When(req.PageID, req.OrdersPerPage).Then(orders, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "ZeroLen",
			args: args{td["ZeroLen"].req, td["ZeroLen"].view},
			prepare: func() {
				data := td["ZeroLen"]
				req := data.req
				orders := data.view

				m.rp.GetRefundsMock.When(req.PageID, req.OrdersPerPage).Then(orders, nil)
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

			got, err := u.GetRefunds(tt.args.req)
			tt.wantErr(t, err)

			require.Equal(t, tt.args.expect, got)
		})
	}
}
