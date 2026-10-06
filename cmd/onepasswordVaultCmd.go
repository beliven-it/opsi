package cmd

import (
	"github.com/spf13/cobra"
)

var onepasswordVaultCmd = &cobra.Command{
	Use:   "vault",
	Args:  cobra.NoArgs,
	Short: "Manage 1password vaults",
	Long:  "Manage 1password vaults",
	Run:   showHelp,
}

func init() {
	onepasswordCmd.AddCommand(onepasswordVaultCmd)
}
