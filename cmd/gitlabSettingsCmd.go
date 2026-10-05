package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabSettingsCmd = &cobra.Command{
	Use:   "settings {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage the settings of Gitlab projects",
	Long:  "Manage the settings of Gitlab projects",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabSettingsCmd)
}
