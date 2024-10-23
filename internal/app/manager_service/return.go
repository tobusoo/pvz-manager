package manager_service

import (
	"context"
	"time"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/metrics"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ManagerService) Return(ctx context.Context, req *desc.ReturnRequest) (*emptypb.Empty, error) {
	const handler = "return"

	timer := time.Now()
	defer func() { metrics.ObserveResponseTime(time.Since(timer), handler) }()

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.ReturnRequest{
		OrderID: req.GetOrderId(),
	}

	err := s.ru.Return(usecase_req)
	if IsServiceError(err) {
		s.sendEvent([]uint64{req.GetOrderId()}, domain.EventOrderGiveCourier, err)
		metrics.IncTotalErrors(handler, err)
	}

	if err == nil {
		metrics.AddTotalReturnedOrders(1, handler)
	}

	return nil, DomainErrToGRPC(err)
}
