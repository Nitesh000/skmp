package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/registry"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the local skill registry cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching latest registry from GitHub...")
		
		idx, err := registry.UpdateCache()
		if err != nil {
			return err
		}

		fmt.Printf("✓ Registry updated! (%d skills, %d bundles)\n", len(idx.Skills), len(idx.Bundles))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
