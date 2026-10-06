package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabProjectCmd = &cobra.Command{
	Use:   "project {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage Gitlab projects",
	Long:  "Manage Gitlab projects",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabProjectCmd)
}
