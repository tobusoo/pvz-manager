package manager_service

import (
	"context"
	"errors"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ManagerService) GiveOrders(ctx context.Context, req *desc.GiveOrdersRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.GiveOrdersRequest{
		Orders: req.GetOrders(),
	}

	err := s.gu.Give(usecase_req)
	err_join := errors.Join(err...)
	s.sendEvent(req.GetOrders(), domain.EventOrderGiveClient, err_join)

	return nil, DomainErrToGRPC(err_join)
}
