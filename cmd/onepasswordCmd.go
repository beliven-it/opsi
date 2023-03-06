package cmd

import (
	"github.com/spf13/cobra"
)

var onepasswordCmd = &cobra.Command{
	Use:   "1password {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "The 1password scope commands",
	Long:  "The 1password scope commands",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(onepasswordCmd)
}
