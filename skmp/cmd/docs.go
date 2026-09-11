package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:    "generate-docs",
	Short:  "Generate man pages for skmp",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		header := &doc.GenManHeader{
			Title:   "SKMP",
			Section: "1",
			Source:  "skmp",
			Manual:  "User Commands",
		}

		err := doc.GenManTree(rootCmd, header, ".")
		if err != nil {
			return err
		}
		fmt.Println("✓ Man pages generated in current directory")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
