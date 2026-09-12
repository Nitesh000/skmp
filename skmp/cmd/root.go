package cmd

import (
	"os"

	"github.com/Nitesh000/skmp/config"
	"github.com/Nitesh000/skmp/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skmp",
	Short: "Skill marketplace for AI agent skills",
	Long:  "A TUI application to install skills from Skill Marketplace with vim bindings for intuitive navigation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.New(config.Version), tea.WithMouseCellMotion(), tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
