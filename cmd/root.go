package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.ozon.dev/chppppr/homework/internal/storage"
)

func init() {
	rootCmd.AddCommand(acceptCmd)
	rootCmd.AddCommand(fulfillCmd)
	rootCmd.AddCommand(returnCmd)
	rootCmd.AddCommand(viewCmd)

	rootCmd.DisableSuggestions = false
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

var (
	st *storage.Storage

	userID         uint64
	orderID        uint64
	expirationDate string

	rootCmd = &cobra.Command{
		Use:  os.Args[0],
		Long: `Utility for managing an order pick-up point`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func SetStorage(s *storage.Storage) {
	st = s
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
