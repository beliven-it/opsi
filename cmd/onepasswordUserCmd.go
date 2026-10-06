package cmd

import (
	"github.com/spf13/cobra"
)

var onepasswordUserCmd = &cobra.Command{
	Use:   "user",
	Args:  cobra.NoArgs,
	Short: "Manage 1password users",
	Long:  "Manage 1password users",
	Run:   showHelp,
}

func init() {
	onepasswordCmd.AddCommand(onepasswordUserCmd)
}
