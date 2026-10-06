package cmd

import (
	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Args:  cobra.NoArgs,
	Short: "Manage Gitlab projects, groups and users",
	Long: `Manage Gitlab projects, groups and users.

Projects, groups and subgroups are referred by their numeric ID, the one
shown on their Gitlab page.`,
	Run: showHelp,
}

func init() {
	rootCmd.AddCommand(gitlabCmd)

	gitlabCmd.AddGroup(
		&cobra.Group{ID: "projects", Title: "Projects:"},
		&cobra.Group{ID: "groups", Title: "Groups and users:"},
	)
	for _, command := range []*cobra.Command{gitlabProjectCmd, gitlabMirrorCmd, gitlabScheduleCmd, gitlabSettingsCmd, gitlabEnvsCmd, gitlabCleanUpPolicyCmd} {
		command.GroupID = "projects"
	}
	for _, command := range []*cobra.Command{gitlabGroupCmd, gitlabSubgroupCmd, gitlabUserCmd} {
		command.GroupID = "groups"
	}
}
