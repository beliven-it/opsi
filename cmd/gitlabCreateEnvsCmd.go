package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// subgroupCmd represents the subgroup command
var gitlabCreateEnvsCmd = &cobra.Command{
	Use:   "envs {project_id} {env_file_path}",
	Args:  cobra.ExactArgs(2),
	Short: "Create ENVs for Gitlab project",
	Long: `
  Create ENVs for a specific Gitlab project.
  You can also create the env for a specific environment using the 
  flag -e. Please see the example section.
	`,
	Example: `	
  Create ENVs for the project 1234.
  opsi gitlab create envs 1234 /file/to/env.yml

  ---

  Create ENVS for the project 1234 but only for staging environment
  opsi gitlab create envs 1234 /file/to/env.yml -e staging
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take project ID
		projectID := args[0]

		// Take env file path
		envFile := args[1]

		// Take optional environement
		env, _ := cmd.Flags().GetString("env")

		// Create environments
		err := gitlab.CreateEnvs(projectID, env, envFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateEnvsCmd)
	gitlabCreateEnvsCmd.Flags().StringP("env", "e", "*", "The environment scope")
}
