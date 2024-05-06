package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabUpdateCleanUpPolicyCmd = &cobra.Command{
	Use:   "cleanup-policy {project_id}",
	Args:  cobra.ExactArgs(1),
	Short: "Update Cleanup Policy for Gitlab project",
	Long: `
  Update Cleanup Policy for a specific Gitlab project.`,
	Example: `	
	Update Cleanup Policy for the project 1234.
  	opsi gitlab update cleanup-policy 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take project ID
		projectID := args[0]

		// Update cleanup policy
		err := gitlab.UpdateCleanUpPolicy(projectID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabUpdateCmd.AddCommand(gitlabUpdateCleanUpPolicyCmd)
}
