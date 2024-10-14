package manager_service

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ManagerService) ViewRefunds(ctx context.Context, req *desc.ViewRefundsRequest) (*desc.ViewRefundsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usecase_req := &dto.ViewRefundsRequest{
		PageID:        req.GetPageId(),
		OrdersPerPage: req.GetOrdersPerPage(),
	}

	orders, err := s.vu.GetRefunds(usecase_req)
	if err != nil {
		return nil, DomainErrToHTPP(err)
	}

	res_orders := OrderViewToProto(orders)
	return &desc.ViewRefundsResponse{Orders: res_orders}, nil
}
