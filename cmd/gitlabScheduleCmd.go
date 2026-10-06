package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Args:  cobra.NoArgs,
	Short: "Manage the default pipeline schedule of Gitlab projects",
	Long:  "Manage the default pipeline schedule of Gitlab projects",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabScheduleCmd)
}
