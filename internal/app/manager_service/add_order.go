package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ManagerService) AddOrder(ctx context.Context, req *desc.AddOrderRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	order := req.GetOrder()
	exp_date := utils.TimeToString(order.GetExpirationDate().AsTime())

	usecase_req := &dto.AddOrderRequest{
		ExpirationDate: exp_date,
		ContainerType:  order.GetPackageType(),
		UserID:         req.GetUserId(),
		OrderID:        req.GetOrderId(),
		Cost:           order.GetCost(),
		Weight:         order.GetWeight(),
		UseTape:        order.GetUseTape(),
	}

	err := s.au.AcceptOrder(usecase_req)
	s.sendEvent([]uint64{req.GetOrderId()}, domain.EventOrderAccepted, err)

	return nil, DomainErrToHTPP(err)
}
