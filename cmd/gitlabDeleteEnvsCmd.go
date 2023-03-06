package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabDeleteEnvsCmd = &cobra.Command{
	Use:   "envs {project_id}",
	Args:  cobra.ExactArgs(1),
	Short: "Delete ENVs for Gitlab project",
	Long: `
  Delete ENVs for a specific Gitlab project.
  You can also delete the env for a specific environment using the 
  flag -e. Please see the example section.
	`,
	Example: `	
  Delete ENVs for the project 1234.
  opsi gitlab delete envs 1234
	
  ---
	
  Delete ENVS for the project 1234 but only for staging environment
  opsi gitlab delete envs 1234 -e staging
 
  ---

  Delete ENVS for the project 1234 without ask for confirmation.
  opsi gitlab delete envs 1234 -f
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take project ID
		projectID := args[0]

		// Take the enviroment env if provided
		env, _ := cmd.Flags().GetString("env")

		// Take the force flag.
		// This can safe your life.
		// TODO: Move force logi here!
		force, _ := cmd.Flags().GetBool("force")

		// Delete environment
		err := gitlab.DeleteEnvs(projectID, env, force)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabDeleteCmd.AddCommand(gitlabDeleteEnvsCmd)
	gitlabDeleteEnvsCmd.Flags().StringP("env", "e", "*", "The environment scope")
	gitlabDeleteEnvsCmd.Flags().BoolP("force", "f", false, "Not ask confirmation to delete")
}
