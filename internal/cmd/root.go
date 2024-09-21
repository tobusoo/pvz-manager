package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
	"gitlab.ozon.dev/chppppr/homework/internal/usecase"
)

func init() {
	rootCmd.AddCommand(acceptCmd)
	rootCmd.AddCommand(giveCmd)
	rootCmd.AddCommand(returnCmd)
	rootCmd.AddCommand(viewCmd)

	rootCmd.DisableSuggestions = false
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if len(os.Args) == 1 {
		rootCmd.Use = ""
	}
}

var (
	st *storage.Storage

	acceptUsecase *usecase.AcceptUsecase
	giveUsecase   *usecase.GiveUsecase
	returnUsecase *usecase.ReturnUsecase
	viewUsecase   *usecase.ViewUsecase

	cost           uint64
	weight         uint64
	pageID         uint64
	userID         uint64
	orderID        uint64
	useTape        bool
	ordersLimit    uint64
	ordersPerPage  uint64
	containerType  string
	expirationDate string
	orders         []uint

	rootCmd = &cobra.Command{
		Use:  "manager",
		Long: `Utility for managing an order pick-up point`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func SetStorage(s *storage.Storage) {
	st = s
	acceptUsecase = usecase.NewAcceptUsecase(st)
	giveUsecase = usecase.NewGiveUsecase(st)
	returnUsecase = usecase.NewReturnUsecase(st)
	viewUsecase = usecase.NewViewUsecase(st)
}

func SetArgs(args []string) {
	rootCmd.SetArgs(args)
}

func Execute() error {
	if st == nil {
		return fmt.Errorf("before Execute() need to set storage")
	}

	return rootCmd.Execute()
}
