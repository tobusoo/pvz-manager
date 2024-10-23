package manager_service

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/metrics"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ManagerService) GiveOrders(ctx context.Context, req *desc.GiveOrdersRequest) (*emptypb.Empty, error) {
	const handler = "give_orders"

	timer := time.Now()
	defer func() { metrics.ObserveResponseTime(time.Since(timer), handler) }()

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.GiveOrdersRequest{
		Orders: req.GetOrders(),
	}

	err := s.gu.Give(usecase_req)
	err_join := errors.Join(err...)
	if IsServiceError(err_join) {
		s.sendEvent(req.GetOrders(), domain.EventOrderGiveClient, err_join)
		metrics.IncTotalErrors(handler, err_join)
	}

	if err == nil {
		metrics.AddTotalIssuedOrders(len(req.GetOrders()), handler)
	}

	return nil, DomainErrToGRPC(err_join)
}
