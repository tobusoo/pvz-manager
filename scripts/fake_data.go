package scripts

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

const DefaultSeed = 123

func init() {
	gofakeit.Seed(DefaultSeed)

	gofakeit.AddFuncLookup("datefuture", gofakeit.Info{
		Category: "custom",
		Output:   "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			startDate := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
			randDate := f.DateRange(startDate.AddDate(0, 6, 0), startDate.AddDate(0, 12, 0))
			return utils.TimeToString(randDate), nil
		},
	})
}

func generateWeightBasedOnContainerType(containerType string) uint64 {
	var weight uint
	switch containerType {
	case "tape":
		weight = gofakeit.UintRange(1, 50000)
	case "package":
		weight = gofakeit.UintRange(1, strategy.PackageMaxWeight)
	case "box":
		weight = gofakeit.UintRange(strategy.PackageMaxWeight, strategy.BoxMaxWeight)
	default:
		weight = gofakeit.UintRange(1, 25200)
	}

	return uint64(weight)
}

func GenerateAddRequests(count int) []*dto.AddOrderRequest {
	res := make([]*dto.AddOrderRequest, 0)
	for i := 0; i < count; i++ {
		var req *dto.AddOrderRequest
		gofakeit.Struct(&req)
		req.Weight = generateWeightBasedOnContainerType(req.ContainerType)
		res = append(res, req)
	}

	return res
}

func GenerateAddRequestsWithUserID(userID uint64, count int) []*dto.AddOrderRequest {
	res := GenerateAddRequests(count)
	for _, req := range res {
		req.UserID = userID
	}

	return res
}

func GenerateGiveRequests(reqs []*dto.AddOrderRequest, limit int) ([]*dto.GiveOrdersRequest, []uint64) {
	limit = min(limit, len(reqs))
	res := make([]*dto.GiveOrdersRequest, 0)
	userIDs := make([]uint64, 0)
	m := make(map[uint64]int)

	for i := 0; i < limit; i++ {
		req := reqs[i]

		id, ok := m[req.UserID]
		if !ok {
			m[req.UserID] = i
			id = i
		}

		if id < i {
			res[id].Orders = append(res[id].Orders, uint(req.OrderID))
		} else {
			res = append(res, &dto.GiveOrdersRequest{Orders: make([]uint, 0)})
			res[i].Orders = append(res[i].Orders, uint(req.OrderID))
			userIDs = append(userIDs, req.UserID)
		}

	}

	return res, userIDs
}

func GenerateRefundRequests(reqs []*dto.GiveOrdersRequest, userIDs []uint64, limit int) []*dto.RefundRequest {
	limit = min(limit, len(reqs))
	res := make([]*dto.RefundRequest, 0)
	for i := 0; i < limit; i++ {
		res = append(res, &dto.RefundRequest{OrderID: uint64(reqs[i].Orders[0]), UserID: userIDs[i]})
	}

	return res
}

func GenerateReturnRequests(reqs []*dto.RefundRequest, limit int) []*dto.ReturnRequest {
	limit = min(limit, len(reqs))
	res := make([]*dto.ReturnRequest, 0)
	for i := 0; i < limit; i++ {
		res = append(res, &dto.ReturnRequest{OrderID: uint64(reqs[i].OrderID)})
	}

	return res
}
