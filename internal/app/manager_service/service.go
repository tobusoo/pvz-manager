package manager_service

import (
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/usecase"
	desc "gitlab.ozon.dev/chppppr/homework/pkg/manager-service/v1"
)

type ManagerService struct {
	st storage.Storage

	au *usecase.AcceptUsecase
	gu *usecase.GiveUsecase
	ru *usecase.ReturnUsecase
	vu *usecase.ViewUsecase

	desc.UnimplementedManagerServiceServer
}

func NewManagerService(st storage.Storage) *ManagerService {
	s := &ManagerService{
		st: st,
		au: usecase.NewAcceptUsecase(st),
		gu: usecase.NewGiveUsecase(st),
		ru: usecase.NewReturnUsecase(st),
		vu: usecase.NewViewUsecase(st),
	}

	return s
}
