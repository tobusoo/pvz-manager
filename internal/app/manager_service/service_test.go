package manager_service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/chppppr/homework/internal/app/manager_service/mock"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestManagerService_AddOrder(t *testing.T) {
	type (
		args struct {
			req *desc.AddOrderRequest
		}

		TestData struct {
			req_dto   *dto.AddOrderRequest
			req_proto *desc.AddOrderRequest
		}
	)

	ctrl := minimock.NewController(t)
	us := mock.NewUsecasesMock(ctrl)
	prod := mock.NewKafkaProducerMock(ctrl)

	mng := NewManagerService(us, us, us, us, prod)
	ctx := context.Background()

	cur_time := utils.CurrentDate()
	time_str := utils.TimeToString(cur_time)
	td := map[string]TestData{
		"Success": {
			req_dto: &dto.AddOrderRequest{
				OrderID:        1,
				UserID:         1,
				ExpirationDate: time_str,
			},
			req_proto: &desc.AddOrderRequest{
				OrderId: 1,
				UserId:  1,
				Order: &desc.Order{
					ExpirationDate: timestamppb.New(cur_time),
				},
			},
		},
		"AlreadyExist": {
			req_dto: &dto.AddOrderRequest{
				OrderID:        2,
				UserID:         2,
				ExpirationDate: time_str,
			},
			req_proto: &desc.AddOrderRequest{
				OrderId: 2,
				UserId:  2,
				Order: &desc.Order{
					ExpirationDate: timestamppb.New(cur_time),
				},
			},
		},
		"ServiceFail": {
			req_dto: &dto.AddOrderRequest{
				OrderID:        3,
				UserID:         3,
				ExpirationDate: time_str,
			},
			req_proto: &desc.AddOrderRequest{
				OrderId: 3,
				UserId:  3,
				Order: &desc.Order{
					ExpirationDate: timestamppb.New(cur_time),
				},
			},
		},
	}

	tests := []struct {
		name    string
		prepare func()
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				req: td["Success"].req_proto,
			},
			prepare: func() {
				data := td["Success"]
				req := data.req_dto

				us.AcceptOrderMock.When(req).Then(nil)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderAccepted, nil, nil).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "AlreadyExist",
			args: args{
				req: td["AlreadyExist"].req_proto,
			},
			prepare: func() {
				data := td["AlreadyExist"]
				req := data.req_dto

				us.AcceptOrderMock.When(req).Then(domain.ErrAlreadyExist)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderAccepted, domain.ErrAlreadyExist, nil).Then(nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "ServiceFail",
			args: args{
				req: td["ServiceFail"].req_proto,
			},
			prepare: func() {
				data := td["ServiceFail"]
				req := data.req_dto

				some_service_error := fmt.Errorf("some bad service error")
				us.AcceptOrderMock.When(req).Then(some_service_error)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderAccepted, nil, some_service_error).Then(nil)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.prepare()

			_, err := mng.AddOrder(ctx, tt.args.req)
			tt.wantErr(t, err)
		})
	}
}

func TestManagerService_GiveOrder(t *testing.T) {
	type (
		args struct {
			req *desc.GiveOrdersRequest
		}

		TestData struct {
			req_dto   *dto.GiveOrdersRequest
			req_proto *desc.GiveOrdersRequest
		}
	)

	ctrl := minimock.NewController(t)
	us := mock.NewUsecasesMock(ctrl)
	prod := mock.NewKafkaProducerMock(ctrl)

	mng := NewManagerService(us, us, us, us, prod)
	ctx := context.Background()

	td := map[string]TestData{
		"Success": {
			req_dto: &dto.GiveOrdersRequest{
				Orders: []uint64{1, 2, 3},
			},
			req_proto: &desc.GiveOrdersRequest{
				Orders: []uint64{1, 2, 3},
			},
		},
		"NotFound": {
			req_dto: &dto.GiveOrdersRequest{
				Orders: []uint64{4},
			},
			req_proto: &desc.GiveOrdersRequest{
				Orders: []uint64{4},
			},
		},
		"ServiceFail": {
			req_dto: &dto.GiveOrdersRequest{
				Orders: []uint64{5},
			},
			req_proto: &desc.GiveOrdersRequest{
				Orders: []uint64{5},
			},
		},
	}

	tests := []struct {
		name    string
		prepare func()
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				req: td["Success"].req_proto,
			},
			prepare: func() {
				data := td["Success"]
				req := data.req_dto

				us.GiveMock.When(req).Then(nil)
				prod.SendMock.When(req.Orders, domain.EventOrderGiveClient, nil, nil).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				req: td["NotFound"].req_proto,
			},
			prepare: func() {
				data := td["NotFound"]
				req := data.req_dto

				err_join := errors.Join(domain.ErrNotFound)
				us.GiveMock.When(req).Then([]error{domain.ErrNotFound})
				prod.SendMock.When(req.Orders, domain.EventOrderGiveClient, err_join, nil).Then(nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "ServiceFail",
			args: args{
				req: td["ServiceFail"].req_proto,
			},
			prepare: func() {
				data := td["ServiceFail"]
				req := data.req_dto

				some_service_error := fmt.Errorf("some bad service error")
				err_join := errors.Join(some_service_error)
				us.GiveMock.When(req).Then([]error{some_service_error})
				prod.SendMock.When(req.Orders, domain.EventOrderGiveClient, nil, err_join).Then(nil)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.prepare()

			_, err := mng.GiveOrders(ctx, tt.args.req)
			tt.wantErr(t, err)
		})
	}
}

func TestManagerService_Refund(t *testing.T) {
	type (
		args struct {
			req *desc.RefundRequest
		}

		TestData struct {
			req_dto   *dto.RefundRequest
			req_proto *desc.RefundRequest
		}
	)

	ctrl := minimock.NewController(t)
	us := mock.NewUsecasesMock(ctrl)
	prod := mock.NewKafkaProducerMock(ctrl)

	mng := NewManagerService(us, us, us, us, prod)
	ctx := context.Background()

	td := map[string]TestData{
		"Success": {
			req_dto: &dto.RefundRequest{
				OrderID: 1,
				UserID:  1,
			},
			req_proto: &desc.RefundRequest{
				OrderId: 1,
				UserId:  1,
			},
		},
		"NotFound": {
			req_dto: &dto.RefundRequest{
				OrderID: 2,
				UserID:  2,
			},
			req_proto: &desc.RefundRequest{
				OrderId: 2,
				UserId:  2,
			},
		},
		"ServiceFail": {
			req_dto: &dto.RefundRequest{
				OrderID: 3,
				UserID:  3,
			},
			req_proto: &desc.RefundRequest{
				OrderId: 3,
				UserId:  3,
			},
		},
	}

	tests := []struct {
		name    string
		prepare func()
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				req: td["Success"].req_proto,
			},
			prepare: func() {
				data := td["Success"]
				req := data.req_dto

				us.AcceptRefundMock.When(req).Then(nil)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderReturned, nil, nil).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				req: td["NotFound"].req_proto,
			},
			prepare: func() {
				data := td["NotFound"]
				req := data.req_dto

				us.AcceptRefundMock.When(req).Then(domain.ErrNotFound)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderReturned, domain.ErrNotFound, nil).Then(nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "ServiceFail",
			args: args{
				req: td["ServiceFail"].req_proto,
			},
			prepare: func() {
				data := td["ServiceFail"]
				req := data.req_dto

				some_service_error := fmt.Errorf("some bad service error")
				us.AcceptRefundMock.When(req).Then(some_service_error)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderReturned, nil, some_service_error).Then(nil)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.prepare()

			_, err := mng.Refund(ctx, tt.args.req)
			tt.wantErr(t, err)
		})
	}
}

func TestManagerService_Return(t *testing.T) {
	type (
		args struct {
			req *desc.ReturnRequest
		}

		TestData struct {
			req_dto   *dto.ReturnRequest
			req_proto *desc.ReturnRequest
		}
	)

	ctrl := minimock.NewController(t)
	us := mock.NewUsecasesMock(ctrl)
	prod := mock.NewKafkaProducerMock(ctrl)

	mng := NewManagerService(us, us, us, us, prod)
	ctx := context.Background()

	td := map[string]TestData{
		"Success": {
			req_dto: &dto.ReturnRequest{
				OrderID: 1,
			},
			req_proto: &desc.ReturnRequest{
				OrderId: 1,
			},
		},
		"NotFound": {
			req_dto: &dto.ReturnRequest{
				OrderID: 2,
			},
			req_proto: &desc.ReturnRequest{
				OrderId: 2,
			},
		},
		"ServiceFail": {
			req_dto: &dto.ReturnRequest{
				OrderID: 3,
			},
			req_proto: &desc.ReturnRequest{
				OrderId: 3,
			},
		},
	}

	tests := []struct {
		name    string
		prepare func()
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				req: td["Success"].req_proto,
			},
			prepare: func() {
				data := td["Success"]
				req := data.req_dto

				us.ReturnMock.When(req).Then(nil)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderGiveCourier, nil, nil).Then(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				req: td["NotFound"].req_proto,
			},
			prepare: func() {
				data := td["NotFound"]
				req := data.req_dto

				us.ReturnMock.When(req).Then(domain.ErrNotFound)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderGiveCourier, domain.ErrNotFound, nil).Then(nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "ServiceFail",
			args: args{
				req: td["ServiceFail"].req_proto,
			},
			prepare: func() {
				data := td["ServiceFail"]
				req := data.req_dto

				some_service_error := fmt.Errorf("some bad service error")
				us.ReturnMock.When(req).Then(some_service_error)
				prod.SendMock.When([]uint64{req.OrderID}, domain.EventOrderGiveCourier, nil, some_service_error).Then(nil)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.prepare()

			_, err := mng.Return(ctx, tt.args.req)
			tt.wantErr(t, err)
		})
	}
}
