package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabSettingsCmd = &cobra.Command{
	Use:   "settings",
	Args:  cobra.NoArgs,
	Short: "Manage the settings of Gitlab projects",
	Long:  "Manage the settings of Gitlab projects",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabSettingsCmd)
}
