package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print skmp version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("skmp v" + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = Version
}
