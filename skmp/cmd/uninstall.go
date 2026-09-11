package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Nitesh000/skmp/harness"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove all skmp data and skills from your system",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Println("This will:")
			fmt.Println("  - Remove all installed skills")
			fmt.Println("  - Remove all symlinks from harness directories")
			fmt.Println("  - Delete ~/.skmp entirely")
			fmt.Println("")
			fmt.Print("Are you sure? (y/N): ")

			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))

			if answer != "y" && answer != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		fmt.Println("Removing skmp data...")
		if err := harness.Uninstall(); err != nil {
			return fmt.Errorf("uninstall: %w", err)
		}

		fmt.Println("✓ All skmp data removed.")
		fmt.Println("")
		fmt.Println("To remove the CLI itself, run:")
		fmt.Println("  npm uninstall -g skmp")
		return nil
	},
}

func init() {
	uninstallCmd.Flags().Bool("force", false, "skip confirmation prompt")
	rootCmd.AddCommand(uninstallCmd)
}
