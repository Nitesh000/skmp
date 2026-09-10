package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/harness"
	"github.com/Nitesh000/skmp/registry"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <skill-name>",
	Short: "Install a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		idx, err := registry.Load()
		if err != nil {
			return fmt.Errorf("registry: %w", err)
		}

		var skill *registry.Skill
		for _, s := range idx.Skills {
			if s.Name == name {
				s := s
				skill = &s
				break
			}
		}

		if skill == nil {
			return fmt.Errorf("skill %q not found", name)
		}

		if harness.IsInstalled(name) {
			fmt.Printf("%s already installed\n", name)
			return nil
		}

		fmt.Printf("installing %s@%s...\n", skill.Name, skill.Version)
		if err := harness.Install(skill.Name, skill.Source); err != nil {
			return err
		}

		fmt.Printf("✓ %s installed\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
