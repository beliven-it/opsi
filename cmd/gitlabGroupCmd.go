package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabGroupCmd = &cobra.Command{
	Use:   "group",
	Args:  cobra.NoArgs,
	Short: "Manage Gitlab groups",
	Long:  "Manage Gitlab groups",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabGroupCmd)
}
