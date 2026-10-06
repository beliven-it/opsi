package cmd

import (
	"fmt"
	"opsi/helpers"
	"os"

	"github.com/spf13/cobra"
)

var gitlabScheduleUpdateCmd = &cobra.Command{
	Use:   "update [project]...",
	Args:  cobra.ArbitraryArgs,
	Short: "Fix the pipeline schedules of Gitlab projects that lost their owner",
	Long: `Fix the active pipeline schedules of Gitlab projects that lost their owner.
Each project can be given as numeric ID or as full path. Without any project,
all the projects are checked.

A schedule runs as its owner, so it stops when the owner is blocked, deactivated
or removed. This command gives those schedules to the user of the token.
The schedules of an active user are left alone, as cron, branch and description.
Nothing is created or activated: projects without an active schedule are ignored.
Use --dry-run to only see the schedules, with their owner.`,
	Example: `
  Show the active schedules of all the projects and who owns them
  opsi gitlab schedule update --dry-run

  ---

  Fix the schedules of more projects, by path or ID
  opsi gitlab schedule update corporate/wiki/beliven-wiki 1234

  ---

  Fix the schedules of all the projects
  opsi gitlab schedule update
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Confirm the action (a dry run changes nothing)
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		if !force && !dryRun {
			helpers.Confirm()
		}

		if err := gitlab.UpdateSchedule(args, dryRun); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabScheduleCmd.AddCommand(gitlabScheduleUpdateCmd)
	gitlabScheduleUpdateCmd.Flags().Bool("dry-run", false, "Only show the schedules and their owner")
	gitlabScheduleUpdateCmd.Flags().BoolP("force", "f", false, "Not ask confirmation")
}
