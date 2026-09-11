package cmd

import (
	"fmt"

	"github.com/Nitesh000/skmp/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print skmp version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("skmp v" + config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = config.Version
}
