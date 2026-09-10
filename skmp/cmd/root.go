package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skmp",
	Short: "Skill marketplace for AI agent skills",
	Long:  "A TUI application to install skills from Skill Marketplace with vim bindings for intuitive navigation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("TUI Comming soon")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
