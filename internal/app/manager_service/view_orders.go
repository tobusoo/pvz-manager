package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ManagerService) ViewOrdersV1(ctx context.Context, req *desc.ViewOrdersRequestV1) (*desc.ViewOrdersResponseV1, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.ViewOrdersRequest{
		UserID:       req.GetUserId(),
		FirstOrderID: req.GetFirstOrderId(),
		OrdersLimit:  req.GetLimit(),
	}

	orders, err := s.vu.GetOrders(usecase_req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res_orders := OrderViewToProto(orders)
	return &desc.ViewOrdersResponseV1{Orders: res_orders}, nil
}
