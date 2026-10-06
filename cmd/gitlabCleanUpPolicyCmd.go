package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabCleanUpPolicyCmd = &cobra.Command{
	Use:   "cleanup-policy",
	Args:  cobra.NoArgs,
	Short: "Manage the cleanup policy of Gitlab projects",
	Long:  "Manage the cleanup policy of Gitlab projects",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabCleanUpPolicyCmd)
}
