package storage_suite

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

type StorageJSONSuite struct {
	suite.Suite
	st *storage.StorageJSON
}

func (s *StorageJSONSuite) SetupSuite() {
	var err error
	ohp := storage.NewOrdersHistory()
	rp := storage.NewRefunds()
	up := storage.NewUsers()

	s.st, err = storage.NewStorage(ohp, rp, up, "test_data/storage_test.json")
	s.Require().NoError(err)
}

func (s *StorageJSONSuite) TestOrderAlreadyExist() {
	userID := uint64(101)
	orderID := uint64(1)
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

func (s *StorageJSONSuite) TestSetOrderStatusNotFound() {
	err := s.st.SetOrderStatus(123, domain.StatusAccepted)
	s.Require().Error(err)
}

func (s *StorageJSONSuite) TestRemoveOrderWrongStatus() {
	err := s.st.CanRemoveOrder(1)
	s.Require().Error(err)
}

func (s *StorageJSONSuite) TestRemoveOrderWrongOrderID() {
	err := s.st.CanRemoveOrder(10101)
	s.Require().Error(err)
}

func (s *StorageJSONSuite) TestGetOrderUserNotFound() {
	_, err := s.st.GetOrder(10, 10)
	s.Require().Error(err)
}

func (s *StorageJSONSuite) TestGetExpirationDateUserNotFound() {
	_, err := s.st.GetExpirationDate(10, 10)
	s.Require().Error(err)
}

func readFromGoldenFile(t *testing.T, path string, actual []byte, isUpdate bool) (want []byte) {
	if isUpdate {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0777)
		require.NoError(t, err)
		defer file.Close()

		_, err = file.Write(actual)
		require.NoError(t, err)

		return actual
	}

	want, err := os.ReadFile(path)
	require.NoError(t, err)
	return
}

func (s *StorageJSONSuite) TestSuccessGetOrders() {
	var expected []domain.OrderView
	actual, err := s.st.GetOrdersByUserID(1, 2, 0)
	s.Require().NoError(err)

	// Для обновления файла актульными данными:
	// actual_bytes, err := json.Marshal(actual)
	// s.Require().NoError(err)
	// expected_bytes := readFromGoldenFile(s.T(), "test_data/GetOrdersValues.json", actual_bytes, true)
	expected_bytes := readFromGoldenFile(s.T(), "test_data/GetOrdersValues.json", nil, false)

	err = json.Unmarshal(expected_bytes, &expected)
	s.Require().NoError(err)
	s.Equal(expected, actual)

	actual, err = s.st.GetOrdersByUserID(1, 3, 0)
	s.Require().NoError(err)
	s.Equal(expected[1:], actual)

	actual, err = s.st.GetOrdersByUserID(1, 3, 3)
	s.Require().NoError(err)
	s.Equal(expected[1:4], actual)
}

func (s *StorageJSONSuite) TestGetOrdersUserNotFound() {
	_, err := s.st.GetOrdersByUserID(101, 2, 0)
	s.Require().Error(err)
}

func (s *StorageJSONSuite) TestGetOrdersWrongFirstUserID() {
	_, err := s.st.GetOrdersByUserID(1, 101, 0)
	s.Require().Error(err)
}

func (s *StorageJSONSuite) TestSuccessGetRefunds() {
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

func (s *StorageJSONSuite) TestRemoveRefundOrderNotFound() {
	err := s.st.RemoveRefund(101)
	s.Require().Error(err)
}

func (s *StorageJSONSuite) TestFailGetRefunds() {
	_, err := s.st.GetRefunds(0, 10)
	s.Require().Error(err)

	_, err = s.st.GetRefunds(1, 0)
	s.Require().Error(err)
}
