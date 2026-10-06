package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabUserCmd = &cobra.Command{
	Use:   "user {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage Gitlab users",
	Long:  "Manage Gitlab users",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabUserCmd)
}
