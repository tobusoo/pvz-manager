package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ManagerService) Refund(ctx context.Context, req *desc.RefundRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.RefundRequest{
		UserID:  req.GetUserId(),
		OrderID: req.GetOrderId(),
	}

	if err := s.au.AcceptRefund(usecase_req); err != nil {
		return nil, DomainErrToHTPP(err)
	}

	return &emptypb.Empty{}, nil
}
