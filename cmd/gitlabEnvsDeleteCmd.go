package cmd

import (
	"opsi/helpers"
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

var gitlabEnvsDeleteCmd = &cobra.Command{
	Use:   "delete {project_id}",
	Args:  cobra.ExactArgs(1),
	Short: "Delete ENVs for Gitlab project",
	Long: `
  Delete ENVs for a specific Gitlab project.
  You can also delete the env for a specific environment using the 
  flag -e. Please see the example section.
	`,
	Example: `	
  Delete ENVs for the project 1234.
  opsi gitlab envs delete 1234
	
  ---
	
  Delete ENVS for the project 1234 but only for staging environment
  opsi gitlab envs delete 1234 -e staging
 
  ---

  Delete ENVS for the project 1234 without ask for confirmation.
  opsi gitlab envs delete 1234 -f
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take project ID
		projectID := args[0]

		// Take the enviroment env if provided
		env, _ := cmd.Flags().GetString("env")

		// Take the force flag.This can safe your life.
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			helpers.Confirm()
		}

		// Delete environment
		err := gitlab.DeleteEnvs(projectID, env)
		if err != nil {
			ui.Fatal(err)
		}
	},
}

func init() {
	gitlabEnvsCmd.AddCommand(gitlabEnvsDeleteCmd)
	gitlabEnvsDeleteCmd.Flags().StringP("env", "e", "*", "The environment scope")
	gitlabEnvsDeleteCmd.Flags().BoolP("force", "f", false, "Not ask confirmation to delete")
}
