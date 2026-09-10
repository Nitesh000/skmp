package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/harness"
	"github.com/Nitesh000/skmp/registry"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <skill> [skill...]",
	Short: "Install one or more skills",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idx, err := registry.Load()
		if err != nil {
			return fmt.Errorf("registry: %w", err)
		}

		// registry by name for fast lookup
		skillMap := map[string]registry.Skill{}
		for _, s := range idx.Skills {
			skillMap[s.Name] = s
		}

		var failed []string
		for _, name := range args {
			s, ok := skillMap[name]

			if !ok {
				fmt.Printf("✗ %s — not found in registry\n", name)
				failed = append(failed, name)
				continue
			}

			if harness.IsInstalled(name) {
				fmt.Printf("  %s — already installed\n", name)
				continue
			}

			fmt.Printf("  installing %s@%s...\n", s.Name, s.Version)
			if err := harness.Install(s.Name, s.Source); err != nil {
				fmt.Printf("✗ %s — %s\n", name, err)
				failed = append(failed, name)
				continue
			}

			fmt.Printf("✓ %s\n", name)
		}

		if len(failed) > 0 {
			return fmt.Errorf("%d skill(s) failed to install", len(failed))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
