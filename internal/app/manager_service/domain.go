package manager_service

import (
	"errors"

	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//gocyclo:ignore
//gocognit:ignore
func DomainErrToHTPP(err error) error {
	if err == nil {
		return err
	}

	if errors.Is(err, domain.ErrWrongInput) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, domain.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, domain.ErrAlreadyExist) {
		return status.Error(codes.AlreadyExists, err.Error())
	} else if errors.Is(err, domain.ErrWrongStatus) ||
		errors.Is(err, domain.ErrExpirationDatePassed) ||
		errors.Is(err, domain.ErrNotExpirationDate) ||
		errors.Is(err, domain.ErrTwoDaysPassed) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

//gocyclo:ignore
//gocognit:ignore
func IsServiceError(err error) bool {
	if errors.Is(err, domain.ErrWrongInput) {
		return false
	} else if errors.Is(err, domain.ErrNotFound) {
		return false
	} else if errors.Is(err, domain.ErrAlreadyExist) {
		return false
	} else if errors.Is(err, domain.ErrWrongStatus) {
		return false
	} else if errors.Is(err, domain.ErrExpirationDatePassed) {
		return false
	} else if errors.Is(err, domain.ErrNotExpirationDate) {
		return false
	} else if errors.Is(err, domain.ErrTwoDaysPassed) {
		return false
	}

	return true
}

func OrderViewToProto(in []domain.OrderView) []*desc.OrderView {
	out := make([]*desc.OrderView, len(in))

	for i, order := range in {
		date, _ := utils.StringToTime(order.ExpirationDate)
		exp_date := timestamppb.New(date)

		proto_order := &desc.Order{
			ExpirationDate: exp_date,
			PackageType:    order.PackageType,
			Cost:           order.Cost,
			Weight:         order.Weight,
			UseTape:        order.UseTape,
		}

		out[i] = &desc.OrderView{
			UserId:  order.UserID,
			OrderId: order.OrderID,
			Order:   proto_order,
		}
	}

	return out
}
