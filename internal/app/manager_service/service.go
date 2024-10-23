package manager_service

import (
	"log"

	"gitlab.ozon.dev/chppppr/homework/internal/clients"
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
		pr clients.KafkaProducer

		desc.UnimplementedManagerServiceServer
	}
)

func NewManagerService(au AcceptUsecase, gu GiveUsecase, ru ReturnUsecase, vu ViewUsecase, pr clients.KafkaProducer) *ManagerService {
	s := &ManagerService{
		au: au,
		gu: gu,
		ru: ru,
		vu: vu,
		pr: pr,
	}

	return s
}

func (s *ManagerService) sendEvent(orderIDs []uint64, event domain.EventType, err error) {
	prod_err := s.pr.Send(orderIDs, event, err)
	if prod_err != nil {
		log.Println("ManagerService.sendEvent() failed: ", prod_err)
	}
}
