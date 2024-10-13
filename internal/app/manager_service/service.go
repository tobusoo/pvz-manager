package manager_service

import (
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
)

type (
	AcceptUsecase interface {
		AcceptOrder(req *dto.AddOrderRequest) error
		AcceptRefund(req *dto.RefundRequest) error
	}

	GiveUsecase interface {
		Give(req *dto.GiveOrdersRequest) []error
	}

	ReturnUsecase interface {
		Return(req *dto.ReturnRequest) error
	}

	ViewUsecase interface {
		GetOrders(req *dto.ViewOrdersRequest) ([]domain.OrderView, error)
		GetRefunds(req *dto.ViewRefundsRequest) ([]domain.OrderView, error)
	}

	Usecases interface {
		AcceptUsecase
		GiveUsecase
		ReturnUsecase
		ViewUsecase
	}

	ManagerService struct {
		au AcceptUsecase
		gu GiveUsecase
		ru ReturnUsecase
		vu ViewUsecase

		desc.UnimplementedManagerServiceServer
	}
)

func NewManagerService(au AcceptUsecase, gu GiveUsecase, ru ReturnUsecase, vu ViewUsecase) *ManagerService {
	s := &ManagerService{
		au: au,
		gu: gu,
		ru: ru,
		vu: vu,
	}

	return s
}
