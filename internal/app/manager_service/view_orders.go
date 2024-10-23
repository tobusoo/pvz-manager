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

func (s *ManagerService) ViewOrders(ctx context.Context, req *desc.ViewOrdersRequest) (*desc.ViewOrdersResponse, error) {
	const handler = "view_orders"

	timer := time.Now()
	defer func() { metrics.ObserveResponseTime(time.Since(timer), handler) }()

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.ViewOrdersRequest{
		UserID:       req.GetUserId(),
		FirstOrderID: req.GetFirstOrderId(),
		OrdersLimit:  req.GetLimit(),
	}

	orders, err := s.vu.GetOrders(usecase_req)
	if IsServiceError(err) {
		metrics.IncTotalErrors(handler, err)
	}

	if err != nil {
		return nil, DomainErrToGRPC(err)
	}

	res_orders := OrderViewToProto(orders)
	return &desc.ViewOrdersResponse{Orders: res_orders}, nil
}
