package manager

import (
	"context"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	manager_service "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ManagerServiceClient struct {
	mng manager_service.ManagerServiceClient
}

func NewManagerServiceClient(mng manager_service.ManagerServiceClient) *ManagerServiceClient {
	return &ManagerServiceClient{mng: mng}
}

func (s *ManagerServiceClient) AddOrder(ctx context.Context, req *dto.AddOrderRequest) error {
	exp_date, err := utils.StringToTime(req.ExpirationDate)
	if err != nil {
		return err
	}

	order := &manager_service.Order{
		ExpirationDate: timestamppb.New(exp_date),
		PackageType:    req.ContainerType,
		UseTape:        req.UseTape,
		Cost:           req.Cost,
		Weight:         req.Weight,
	}

	req_proto := &manager_service.AddOrderRequest{
		OrderId: req.OrderID,
		UserId:  req.UserID,
		Order:   order,
	}

	_, err = s.mng.AddOrder(ctx, req_proto)
	return err
}

func (s *ManagerServiceClient) Refund(ctx context.Context, req *dto.RefundRequest) error {
	req_proto := &manager_service.RefundRequest{
		UserId:  req.UserID,
		OrderId: req.OrderID,
	}

	_, err := s.mng.Refund(ctx, req_proto)
	return err
}

func (s *ManagerServiceClient) GiveOrders(ctx context.Context, req *dto.GiveOrdersRequest) error {
	req_proto := &manager_service.GiveOrdersRequest{
		Orders: req.Orders,
	}

	_, err := s.mng.GiveOrders(ctx, req_proto)
	return err
}

func (s *ManagerServiceClient) Return(ctx context.Context, req *dto.ReturnRequest) error {
	req_proto := &manager_service.ReturnRequest{
		OrderId: req.OrderID,
	}

	_, err := s.mng.Return(ctx, req_proto)
	return err
}

func (s *ManagerServiceClient) ViewOrders(ctx context.Context, req *dto.ViewOrdersRequest) (*dto.ViewOrdersResponse, error) {
	req_proto := &manager_service.ViewOrdersRequest{
		UserId:       req.UserID,
		FirstOrderId: req.FirstOrderID,
		Limit:        req.OrdersLimit,
	}

	res_proto, err := s.mng.ViewOrders(ctx, req_proto)
	res := orderViewToDomain(res_proto.GetOrders())

	return &dto.ViewOrdersResponse{Orders: res}, err
}

func (s *ManagerServiceClient) ViewRefunds(ctx context.Context, req *dto.ViewRefundsRequest) (*dto.ViewRefundsResponse, error) {
	req_proto := &manager_service.ViewRefundsRequest{
		PageId:        req.PageID,
		OrdersPerPage: req.OrdersPerPage,
	}

	res_proto, err := s.mng.ViewRefunds(ctx, req_proto)
	res := orderViewToDomain(res_proto.GetOrders())

	return &dto.ViewRefundsResponse{Orders: res}, err
}

func orderViewToDomain(in []*manager_service.OrderView) []domain.OrderView {
	out := make([]domain.OrderView, len(in))

	for i, orderView := range in {
		order := orderView.GetOrder()
		date_str := utils.TimeToString(order.GetExpirationDate().AsTime())

		order_view_domain := &domain.Order{
			ExpirationDate: date_str,
			PackageType:    order.GetPackageType(),
			Cost:           order.GetCost(),
			Weight:         order.GetWeight(),
			UseTape:        order.GetUseTape(),
		}

		out[i] = domain.OrderView{
			UserID:  orderView.GetUserId(),
			OrderID: orderView.GetOrderId(),
			Order:   order_view_domain,
		}
	}

	return out
}
