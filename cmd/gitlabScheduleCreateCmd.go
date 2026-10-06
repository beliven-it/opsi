package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabScheduleCreateCmd = &cobra.Command{
	Use:   "create {project}...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create the default pipeline schedule for existing Gitlab projects",
	Long: `Create the default pipeline schedule for one or more existing Gitlab projects.
Each project can be given as numeric ID or as full path.

The schedule is weekly, on a working day between 09:00 and 18:00 (Europe/Rome).
The day and the time are chosen among the ones used by the fewest schedules,
so the pipelines are spread instead of starting all together.
It runs on the "main" branch when the project has it, otherwise on the default branch.
Projects that already have an active schedule are left alone.`,
	Example: `
  Create the schedule for a project
  opsi gitlab schedule create corporate/wiki/beliven-wiki

  ---

  Create the schedule for more projects, by path or ID
  opsi gitlab schedule create corporate/wiki/beliven-wiki 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gitlab.CreateSchedule(args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabScheduleCmd.AddCommand(gitlabScheduleCreateCmd)
}
