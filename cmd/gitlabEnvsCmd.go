package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabEnvsCmd = &cobra.Command{
	Use:   "envs",
	Args:  cobra.NoArgs,
	Short: "Manage the environment variables of a Gitlab project",
	Long:  "Manage the environment variables of a Gitlab project",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabEnvsCmd)
}
