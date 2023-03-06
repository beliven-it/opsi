package cmd

import (
	"github.com/spf13/cobra"
)

var postmarkListCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.ExactArgs(1),
	Short: "Allow to list entities",
	Long:  "Allow to list entities",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	postmarkCmd.AddCommand(postmarkListCmd)
}
