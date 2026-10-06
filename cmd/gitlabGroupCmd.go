package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabGroupCmd = &cobra.Command{
	Use:   "group {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage Gitlab groups",
	Long:  "Manage Gitlab groups",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabGroupCmd)
}
