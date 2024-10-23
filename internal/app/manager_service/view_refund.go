package manager_service

import (
	"context"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/metrics"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ManagerService) ViewRefunds(ctx context.Context, req *desc.ViewRefundsRequest) (*desc.ViewRefundsResponse, error) {
	const handler = "view_refunds"

	timer := time.Now()
	defer func() { metrics.ObserveResponseTime(time.Since(timer), handler) }()

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.ViewRefundsRequest{
		PageID:        req.GetPageId(),
		OrdersPerPage: req.GetOrdersPerPage(),
	}

	orders, err := s.vu.GetRefunds(usecase_req)
	if IsServiceError(err) {
		metrics.IncTotalErrors(handler, err)
	}

	if err != nil {
		return nil, DomainErrToGRPC(err)
	}

	res_orders := OrderViewToProto(orders)
	return &desc.ViewRefundsResponse{Orders: res_orders}, nil
}
