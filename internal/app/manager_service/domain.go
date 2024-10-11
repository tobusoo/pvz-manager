package manager_service

import (
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrderViewToProto(in []domain.OrderView) []*desc.OrderViewV1 {
	out := make([]*desc.OrderViewV1, 0, len(in))

	for _, order := range in {
		date, _ := utils.StringToTime(order.ExpirationDate)
		exp_date := timestamppb.New(date)

		proto_order := &desc.OrderV1{
			ExpirationDate: exp_date,
			PackageType:    order.PackageType,
			Cost:           order.Cost,
			Weight:         order.Weight,
			UseTape:        order.UseTape,
		}

		res := &desc.OrderViewV1{
			UserId:  order.UserID,
			OrderId: order.OrderID,
			Order:   proto_order,
		}

		out = append(out, res)
	}

	return out
}
