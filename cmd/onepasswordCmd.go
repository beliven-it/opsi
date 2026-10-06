package cmd

import (
	"github.com/spf13/cobra"
)

var onepasswordCmd = &cobra.Command{
	Use:   "1password",
	Args:  cobra.NoArgs,
	Short: "Manage 1password users and vaults",
	Long: `Manage 1password users and vaults.
It needs the 1password CLI (op) installed and signed in.`,
	Run: showHelp,
}

func init() {
	rootCmd.AddCommand(onepasswordCmd)
}
