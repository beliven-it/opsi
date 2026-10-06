package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabScheduleCmd = &cobra.Command{
	Use:   "schedule {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage the default pipeline schedule of Gitlab projects",
	Long:  "Manage the default pipeline schedule of Gitlab projects",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabScheduleCmd)
}
