package benchmark

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/chppppr/homework/internal/domain"
	"gitlab.ozon.dev/chppppr/homework/internal/domain/strategy"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/utils"
)

const StoragePath = "storage_bench.json"

func newStorage() (*storage.StorageJSON, error) {
	ohp := storage.NewOrdersHistory()
	rp := storage.NewRefunds()
	up := storage.NewUsers()

	return storage.NewStorage(ohp, rp, up, StoragePath)
}

func BenchmarkAddOrder(b *testing.B) {
	benches := []struct {
		name   string
		orders int
	}{
		{
			name:   "100",
			orders: 100,
		},
		{
			name:   "1000",
			orders: 1000,
		},
		{
			name:   "10000",
			orders: 10000,
		},
		{
			name:   "100000",
			orders: 100000,
		},
		{
			name:   "1000000",
			orders: 1000000,
		},
	}

	cs := strategy.ContainerTypeMap["package"]
	cs.UseTape()
	order, err := domain.NewOrder(100, 100, utils.CurrentDateString(), cs)
	require.NoError(b, err)

	b.ResetTimer()
	for _, bc := range benches {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				os.Remove(StoragePath)
				st, err := newStorage()
				require.NoError(b, err)
				b.StartTimer()

				for i := 0; i < bc.orders; i++ {
					st.AddOrder(uint64(i), uint64(i), order)
				}
			}
		})
	}

	os.Remove(StoragePath)
}

func BenchmarkAddOrderSingleUser(b *testing.B) {
	benches := []struct {
		name   string
		orders int
	}{
		{
			name:   "100",
			orders: 100,
		},
		{
			name:   "1000",
			orders: 1000,
		},
		{
			name:   "10000",
			orders: 10000,
		},
		{
			name:   "100000",
			orders: 100000,
		},
		{
			name:   "1000000",
			orders: 1000000,
		},
	}

	cs := strategy.ContainerTypeMap[""]
	order, err := domain.NewOrder(100, 100, utils.CurrentDateString(), cs)
	require.NoError(b, err)

	b.ResetTimer()
	for _, bc := range benches {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				os.Remove(StoragePath)
				st, err := newStorage()
				require.NoError(b, err)
				b.StartTimer()

				for i := 0; i < bc.orders; i++ {
					st.AddOrder(1, uint64(i), order)
				}
			}
		})
	}

	os.Remove(StoragePath)
}
