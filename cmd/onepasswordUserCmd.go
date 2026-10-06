package cmd

import (
	"github.com/spf13/cobra"
)

var onepasswordUserCmd = &cobra.Command{
	Use:   "user {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage 1password users",
	Long:  "Manage 1password users",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	onepasswordCmd.AddCommand(onepasswordUserCmd)
}
