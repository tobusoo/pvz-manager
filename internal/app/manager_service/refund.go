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

func (s *ManagerService) Refund(ctx context.Context, req *desc.RefundRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.RefundRequest{
		UserID:  req.GetUserId(),
		OrderID: req.GetOrderId(),
	}

	err := s.au.AcceptRefund(usecase_req)
	s.sendEvent([]uint64{req.GetOrderId()}, domain.EventOrderReturned, err)

	return nil, DomainErrToHTPP(err)
}
