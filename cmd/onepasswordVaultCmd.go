package cmd

import (
	"github.com/spf13/cobra"
)

var onepasswordVaultCmd = &cobra.Command{
	Use:   "vault {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage 1password vaults",
	Long:  "Manage 1password vaults",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	onepasswordCmd.AddCommand(onepasswordVaultCmd)
}
