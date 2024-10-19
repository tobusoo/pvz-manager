package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ManagerService) Return(ctx context.Context, req *desc.ReturnRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.ReturnRequest{
		OrderID: req.GetOrderId(),
	}

	err := s.ru.Return(usecase_req)
	s.sendEvent([]uint64{req.GetOrderId()}, domain.EventOrderGiveCourier, err)

	return nil, DomainErrToHTPP(err)
}
