package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
	storage_suite "gitlab.ozon.dev/chppppr/homework/tests/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, &storage_suite.StorageSuite{})
}

func TestStorageSuccessAddOrder(t *testing.T) {
	t.Parallel()

	ohp := storage.NewOrdersHistory()
	rp := storage.NewRefunds()
	up := storage.NewUsers()
	path := "storage_TestStorageSuccessAdd.json"
	defer os.Remove(path)

	st, err := storage.NewStorage(ohp, rp, up, path)
	require.NoError(t, err)

	userID := uint64(1)
	orderID := uint64(1)
	cost := uint64(100)
	weight := uint64(80)
	cs := strategy.ContainerTypeMap["box"]
	cs.UseTape()

	expect_order := &domain.Order{
		ExpirationDate: utils.CurrentDateString(),
		Cost:           cost + strategy.CostBox + strategy.CostTape,
		Weight:         weight,
		PackageType:    "taped box",
		UseTape:        true,
	}

	order, err := domain.NewOrder(cost, weight, utils.CurrentDateString(), cs)
	require.NoError(t, err)
	require.Equal(t, expect_order, order)

	err = st.AddOrder(userID, orderID, order)
	require.NoError(t, err)

	order, err = st.GetOrder(userID, orderID)
	require.NoError(t, err)
	require.Equal(t, expect_order, order)

	expDate, err := st.GetExpirationDate(userID, orderID)
	require.NoError(t, err)
	require.Equal(t, utils.CurrentDate(), expDate)
}

func TestStorageSuccessAddRefund(t *testing.T) {
	t.Parallel()

	ohp := storage.NewOrdersHistory()
	rp := storage.NewRefunds()
	up := storage.NewUsers()
	path := "storage_TestStorageSuccessAddRefund.json"
	defer os.Remove(path)

	st, err := storage.NewStorage(ohp, rp, up, path)
	require.NoError(t, err)

	userID := uint64(1)
	orderID := uint64(1)
	order := &domain.Order{
		ExpirationDate: utils.CurrentDateString(),
		Cost:           123,
		Weight:         80,
		PackageType:    "default",
		UseTape:        false,
	}

	err = st.AddRefund(userID, orderID, order)
	require.NoError(t, err)
}

func TestStorageSuccessRemoveOrder(t *testing.T) {
	t.Parallel()

	ohp := storage.NewOrdersHistory()
	rp := storage.NewRefunds()
	up := storage.NewUsers()
	path := "storage_TestStorageSuccessRemoveOrder.json"
	defer os.Remove(path)

	st, err := storage.NewStorage(ohp, rp, up, path)
	require.NoError(t, err)

	userID := uint64(1)
	orderID := uint64(1)
	cost := uint64(100)
	weight := uint64(80)
	cs := strategy.ContainerTypeMap["box"]
	cs.UseTape()

	order, err := domain.NewOrder(cost, weight, utils.CurrentDateString(), cs)
	require.NoError(t, err)

	give_orders := []uint64{orderID, orderID + 1, orderID + 2}

	for id := range give_orders {
		err = st.AddOrder(userID, uint64(id), order)
		require.NoError(t, err)
	}

	for id := range give_orders {
		err = st.CanRemoveOrder(uint64(id))
		require.NoError(t, err)

		err = st.RemoveOrder(uint64(id), domain.StatusGiveClient)
		require.NoError(t, err)
	}
}

func TestStorageSuccessReturn(t *testing.T) {
	t.Parallel()

	ohp := storage.NewOrdersHistory()
	rp := storage.NewRefunds()
	up := storage.NewUsers()
	path := "storage_TestStorageSuccessReturn.json"
	defer os.Remove(path)

	st, err := storage.NewStorage(ohp, rp, up, path)
	require.NoError(t, err)

	userID := uint64(1)
	orderID := uint64(1)
	cost := uint64(100)
	weight := uint64(80)
	cs := strategy.ContainerTypeMap["box"]
	cs.UseTape()

	order, err := domain.NewOrder(cost, weight, utils.CurrentDateString(), cs)
	require.NoError(t, err)

	// Заказ принят
	err = st.AddOrder(userID, orderID, order)
	require.NoError(t, err)

	// Заказ забрали
	err = st.CanRemoveOrder(orderID)
	require.NoError(t, err)
	err = st.RemoveOrder(orderID, domain.StatusGiveClient)
	require.NoError(t, err)

	// Заказ вернули на ПВЗ
	err = st.AddRefund(userID, orderID, order)
	require.NoError(t, err)

	// Заказ вернули курьеру
	err = st.RemoveRefund(orderID)
	require.NoError(t, err)
}
