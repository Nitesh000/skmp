package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/harness"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills across all harness",
	RunE: func(cmd *cobra.Command, args []string) error {
		linked, err := harness.Sync()
		if err != nil {
			return err
		}

		if linked == 0 {
			fmt.Printf("everything already in sync")
			return nil
		}
		fmt.Printf("✓ %d link(s) created\n", linked)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
