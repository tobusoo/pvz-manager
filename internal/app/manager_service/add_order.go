package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ManagerService) AddOrderV1(ctx context.Context, req *desc.AddOrderRequestV1) (*desc.AddOrderResponseV1, error) {
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

	if err := s.au.AcceptOrder(usecase_req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.AddOrderResponseV1{}, nil
}
