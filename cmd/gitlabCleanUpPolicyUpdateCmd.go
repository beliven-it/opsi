package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabCleanUpPolicyUpdateCmd = &cobra.Command{
	Use:   "update {project_id}",
	Args:  cobra.MaximumNArgs(1),
	Short: "Update Cleanup Policy for Gitlab project",
	Long: `
  Update Cleanup Policy for a specific Gitlab project.`,
	Example: `	
	Update Cleanup Policy for the project 1234.
  	opsi gitlab cleanup-policy update 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		projectID := ""
		if len(args) > 0 {
			// Take project ID
			projectID = args[0]
		}

		// Update cleanup policy
		err := gitlab.UpdateCleanUpPolicy(projectID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabCleanUpPolicyCmd.AddCommand(gitlabCleanUpPolicyUpdateCmd)
}
