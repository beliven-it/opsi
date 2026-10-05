package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabEnvsCmd = &cobra.Command{
	Use:   "envs {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage the environment variables of a Gitlab project",
	Long:  "Manage the environment variables of a Gitlab project",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabEnvsCmd)
}
