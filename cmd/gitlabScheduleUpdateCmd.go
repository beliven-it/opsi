package cmd

import (
	"fmt"
	"opsi/helpers"
	"os"

	"github.com/spf13/cobra"
)

var gitlabScheduleUpdateCmd = &cobra.Command{
	Use:   "update {project}...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Take the ownership of the pipeline schedules of Gitlab projects",
	Long: `Take the ownership of the pipeline schedules of one or more Gitlab projects.
Each project can be given as numeric ID or as full path.

A schedule runs as its owner, so it stops when the owner is blocked, deactivated
or removed. This command gives the schedules to the user of the token, if they
belong to someone else. Cron, branch and description are not touched.
Use --dry-run to only see the schedules, with their owner.`,
	Example: `
  Show the schedules of a project and who owns them
  opsi gitlab schedule update corporate/wiki/beliven-wiki --dry-run

  ---

  Take the ownership of the schedules of more projects, by path or ID
  opsi gitlab schedule update corporate/wiki/beliven-wiki 1234
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
