package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ManagerService) ReturnV1(ctx context.Context, req *desc.ReturnRequestV1) (*desc.ReturnResponseV1, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.ReturnRequest{
		OrderID: req.GetOrderId(),
	}

	if err := s.ru.Return(usecase_req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.ReturnResponseV1{}, nil
}
