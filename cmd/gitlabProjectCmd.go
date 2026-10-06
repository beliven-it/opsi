package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabProjectCmd = &cobra.Command{
	Use:   "project",
	Args:  cobra.NoArgs,
	Short: "Manage Gitlab projects",
	Long:  "Manage Gitlab projects",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabProjectCmd)
}
