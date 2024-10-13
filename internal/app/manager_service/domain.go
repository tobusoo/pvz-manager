package manager_service

import (
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrderViewToProto(in []domain.OrderView) []*desc.OrderView {
	out := make([]*desc.OrderView, 0, len(in))

	for _, order := range in {
		date, _ := utils.StringToTime(order.ExpirationDate)
		exp_date := timestamppb.New(date)

		proto_order := &desc.Order{
			ExpirationDate: exp_date,
			PackageType:    order.PackageType,
			Cost:           order.Cost,
			Weight:         order.Weight,
			UseTape:        order.UseTape,
		}

		res := &desc.OrderView{
			UserId:  order.UserID,
			OrderId: order.OrderID,
			Order:   proto_order,
		}

		out = append(out, res)
	}

	return out
}
