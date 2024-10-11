package manager_service

import (
	"context"
	"errors"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ManagerService) GiveOrdersV1(ctx context.Context, req *desc.GiveOrdersRequestV1) (*desc.GiveOrdersResponseV1, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.GiveOrdersRequest{
		Orders: req.GetOrders(),
	}

	if err := s.gu.Give(usecase_req); err != nil {
		err_join := errors.Join(err...)
		return nil, status.Error(codes.Internal, err_join.Error())
	}

	return &desc.GiveOrdersResponseV1{}, nil
}
