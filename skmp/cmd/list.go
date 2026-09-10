package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/harness"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := harness.InstalledSkills()
		if err != nil {
			return err
		}

		if len(names) == 0 {
			fmt.Println("no skills installed")
			return nil
		}

		fmt.Printf("installed (%d):\n", len(names))
		for _, n := range names {
			fmt.Printf("  ✓ %s\n", n)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
