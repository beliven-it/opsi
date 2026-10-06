package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabCleanUpPolicyCmd = &cobra.Command{
	Use:   "cleanup-policy {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage the cleanup policy of Gitlab projects",
	Long:  "Manage the cleanup policy of Gitlab projects",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabCleanUpPolicyCmd)
}
