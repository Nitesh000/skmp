package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/harness"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <skill-name>",
	Short: "Uninstall a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if !harness.IsInstalled(name) {
			return fmt.Errorf("skill %q not installed", name)
		}

		if err := harness.Remove(name); err != nil {
			return err
		}

		fmt.Printf("✓ %s removed\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
