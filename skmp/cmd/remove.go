package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/harness"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <skill> [skill...]",
	Short: "Uninstall one or more skills",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var failed []string
		for _, name := range args {
			if !harness.IsInstalled(name) {
				fmt.Printf("✗ %s — not installed\n", name)
				failed = append(failed, name)
				continue
			}

			if err := harness.Remove(name); err != nil {
				fmt.Printf("✗ %s — %s\n", name, err)
				failed = append(failed, name)
				continue
			}
			fmt.Printf("✓ %s removed\n", name)
		}

		if len(failed) > 0 {
			return fmt.Errorf("%d skill(s) falied to remove", len(failed))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
