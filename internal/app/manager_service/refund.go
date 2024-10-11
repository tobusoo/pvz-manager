package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ManagerService) RefundV1(ctx context.Context, req *desc.RefundRequestV1) (*desc.RefundResponseV1, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.RefundRequest{
		UserID:  req.GetUserId(),
		OrderID: req.GetOrderId(),
	}

	if err := s.au.AcceptRefund(usecase_req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.RefundResponseV1{}, nil
}
