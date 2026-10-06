package cmd

import (
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

var gitlabScheduleCreateCmd = &cobra.Command{
	Use:   "create {project_id}...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create the default pipeline schedule for existing Gitlab projects",
	Long: `Create the default pipeline schedule for one or more existing Gitlab projects.
Each project is given by its numeric ID.

The schedule is weekly, on a working day between 09:00 and 18:00 (Europe/Rome).
The day and the time are chosen among the ones used by the fewest schedules,
so the pipelines are spread instead of starting all together.
It runs on the "main" branch when the project has it, otherwise on the default branch.
Projects that already have an active schedule are left alone.`,
	Example: `
  Create the schedule for a project
  opsi gitlab schedule create 239

  ---

  Create the schedule for more projects
  opsi gitlab schedule create 239 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := projectIDs(args)
		if err != nil {
			ui.Fatal(err)
		}

		if err := gitlab.CreateSchedule(ids); err != nil {
			ui.Fatal(err)
		}
	},
}

func init() {
	gitlabScheduleCmd.AddCommand(gitlabScheduleCreateCmd)
}
