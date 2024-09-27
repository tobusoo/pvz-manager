package storage_suite

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
	"gitlab.ozon.dev/chppppr/homework/internal/dto"
	"gitlab.ozon.dev/chppppr/homework/internal/storage/postgres"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	"gitlab.ozon.dev/chppppr/homework/scripts"
)

const (
	TotalAddRequests    = 10000
	TotalGiveRequests   = 5000
	TotalRefundRequsts  = 2500
	TotalReturnRequests = 1250
)

func init() {
	_ = godotenv.Load()
}

type StorageDBSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	st   *postgres.StorageDB
}

func (s *StorageDBSuite) SetupSuite() {
	var err error
	psqlDSN, ok := os.LookupEnv("POSTGRESQL_TEST_DSN")
	if !ok {
		s.FailNow("Not found POSTGRESQL_TEST_DSN at .env")
	}

	ctx := context.Background()
	s.pool, err = pgxpool.Connect(ctx, psqlDSN)
	s.Require().NoError(err)

	txManager := postgres.NewTxManager(s.pool)
	pgPepo := postgres.NewRepoPG(txManager)
	s.st = postgres.NewStorageDB(ctx, txManager, pgPepo)

	s.generateFakeData()
}

func (s *StorageDBSuite) TearDownSuite() {
	s.pool.Close()
}

func (s *StorageDBSuite) execAddRequests() []*dto.AddOrderRequest {
	bad_requests := 0

	log.Printf("Generating %d AddRequests\n", TotalAddRequests)
	addRequests := scripts.GenerateAddRequests(TotalAddRequests - 50)
	for _, req := range addRequests {
		cs := strategy.ContainerTypeMap[req.ContainerType]
		if req.UseTape && req.ContainerType != "tape" {
			cs.UseTape()
		} else {
			req.UseTape = false
		}

		order, _ := domain.NewOrder(req.Cost, req.Weight, req.ExpirationDate, cs)

		if err := s.st.AddOrder(req.UserID, req.OrderID, order); err != nil {
			bad_requests++
		}
	}

	addRequestsUserID := scripts.GenerateAddRequestsWithUserID(12345678, 50)
	for _, req := range addRequestsUserID {
		cs := strategy.ContainerTypeMap[req.ContainerType]
		if req.UseTape && req.ContainerType != "tape" {
			cs.UseTape()
		} else {
			req.UseTape = false
		}

		order, _ := domain.NewOrder(req.Cost, req.Weight, req.ExpirationDate, cs)

		if err := s.st.AddOrder(req.UserID, req.OrderID, order); err != nil {
			bad_requests++
		}
	}

	log.Printf("Generated and exec %d AddRequests: %d bad requests\n", TotalAddRequests, bad_requests)
	return addRequests
}

func (s *StorageDBSuite) execGiveRequests(addRequests []*dto.AddOrderRequest) ([]*dto.GiveOrdersRequest, []uint64) {
	bad_requests := 0

	log.Printf("Generating %d GiveRequests\n", TotalGiveRequests)
	giveRequests, userIDs := scripts.GenerateGiveRequests(addRequests, TotalGiveRequests)
	s.Equal(len(giveRequests), len(userIDs))

	for i := 0; i < len(giveRequests); i++ {
		orders := make([]uint64, 0, len(giveRequests[i].Orders))
		for _, v := range giveRequests[i].Orders {
			orders = append(orders, uint64(v))
		}

		if err := s.st.RemoveOrders(orders, domain.StatusGiveClient); err != nil {
			bad_requests++
		}
	}

	log.Printf("Generated and exec %d GiveRequests: %d bad requests\n", TotalGiveRequests, bad_requests)
	return giveRequests, userIDs
}

func (s *StorageDBSuite) execRefundRequests(giveRequests []*dto.GiveOrdersRequest, userIDs []uint64) []*dto.RefundRequest {
	bad_requests := 0

	log.Printf("Generating %d RefundRequests\n", TotalRefundRequsts)
	refundRequests := scripts.GenerateRefundRequests(giveRequests, userIDs, TotalRefundRequsts)
	for _, req := range refundRequests {
		if err := s.st.AddRefund(req.UserID, req.OrderID, nil); err != nil {
			bad_requests++
		}
	}

	log.Printf("Generated and exec %d RefundRequests: %d bad requests\n", TotalRefundRequsts, bad_requests)
	return refundRequests
}

func (s *StorageDBSuite) execReturnRequests(refundRequests []*dto.RefundRequest) {
	bad_requests := 0

	log.Printf("Generating %d ReturnRequests\n", TotalReturnRequests)
	returnRequests := scripts.GenerateReturnRequests(refundRequests, TotalReturnRequests)
	for _, req := range returnRequests {
		if err := s.st.RemoveRefund(req.OrderID); err != nil {
			bad_requests++
		}
	}

	log.Printf("Generated and exec %d ReturnRequests: %d bad requests\n", TotalReturnRequests, bad_requests)
}

func (s *StorageDBSuite) generateFakeData() {
	addRequests := s.execAddRequests()
	giveRequests, userIDs := s.execGiveRequests(addRequests)
	refundRequests := s.execRefundRequests(giveRequests, userIDs)
	s.execReturnRequests(refundRequests)
}

func (s *StorageDBSuite) TestOrderAlreadyExist() {
	userID := uint64(12345678)
	orderID := uint64(221482238527448200)
	cost := uint64(100)
	weight := uint64(80)
	cs := strategy.ContainerTypeMap[""]
	cs.UseTape()

	order, err := domain.NewOrder(cost, weight, utils.CurrentDateString(), cs)
	s.Require().NoError(err)

	err = s.st.AddOrder(userID, orderID, order)
	s.Require().Error(err)

	err = s.st.AddOrderStatus(orderID, userID, domain.StatusAccepted, order)
	s.Require().Error(err)
}

func (s *StorageDBSuite) TestSetOrderStatusNotFound() {
	err := s.st.SetOrderStatus(123, domain.StatusAccepted)
	s.Require().Error(err)
}

func (s *StorageDBSuite) TestRemoveOrderWrongStatus() {
	err := s.st.CanRemoveOrder(2434254602525190348)
	s.Require().Error(err)
}

func (s *StorageDBSuite) TestRemoveOrderWrongOrderID() {
	err := s.st.CanRemoveOrder(123)
	s.Require().Error(err)
}

func (s *StorageDBSuite) TestGetOrderUserNotFound() {
	_, err := s.st.GetOrder(10, 10)
	s.Require().Error(err)
}

func (s *StorageDBSuite) TestGetExpirationDateUserNotFound() {
	_, err := s.st.GetExpirationDate(10, 10)
	s.Require().Error(err)
}

func (s *StorageDBSuite) TestSuccessGetOrders() {
	var expected []domain.OrderView
	userID := uint64(12345678)
	actual, err := s.st.GetOrdersByUserID(userID, 0, 10)
	s.Require().NoError(err)

	// Для обновления файла актульными данными:
	// actual_bytes, err := json.Marshal(actual)
	// s.Require().NoError(err)
	// expected_bytes := readFromGoldenFile(s.T(), "test_data/GetOrdersValues.json", actual_bytes, true)
	expected_bytes := readFromGoldenFile(s.T(), "test_data/GetOrdersValues.json", nil, false)

	err = json.Unmarshal(expected_bytes, &expected)
	s.Require().NoError(err)
	s.Equal(expected, actual)

	actual, err = s.st.GetOrdersByUserID(userID, expected[1].OrderID, 9)
	s.Require().NoError(err)
	s.Equal(expected[1:], actual)

	actual, err = s.st.GetOrdersByUserID(userID, expected[1].OrderID, 3)
	s.Require().NoError(err)
	s.Equal(expected[1:4], actual)
}

func (s *StorageDBSuite) TestSuccessGetRefunds() {
	var expected []domain.OrderView
	actual, err := s.st.GetRefunds(1, 10)
	s.Require().NoError(err)

	// Для обновления файла актульными данными:
	// actual_bytes, err := json.Marshal(actual)
	// s.Require().NoError(err)
	// expected_bytes := readFromGoldenFile(s.T(), "test_data/GetRefundsValues.json", actual_bytes, true)
	expected_bytes := readFromGoldenFile(s.T(), "test_data/GetRefundsValues.json", nil, false)

	err = json.Unmarshal(expected_bytes, &expected)
	s.Require().NoError(err)
	s.Equal(expected, actual)

	actual, err = s.st.GetRefunds(2, 1)
	s.Require().NoError(err)
	s.Equal(expected[1:2], actual)

	actual, err = s.st.GetRefunds(1, 4)
	s.Require().NoError(err)
	s.Equal(expected[0:4], actual)

	actual, err = s.st.GetRefunds(3, 1)
	s.Require().NoError(err)
	s.Equal(expected[2:3], actual)
}

func (s *StorageDBSuite) TestRemoveRefundOrderNotFound() {
	err := s.st.RemoveRefund(101)
	s.Require().Error(err)
}

func (s *StorageDBSuite) TestFailGetRefunds() {
	_, err := s.st.GetRefunds(0, 10)
	s.Require().Error(err)
}
