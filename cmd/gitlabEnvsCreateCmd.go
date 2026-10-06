package cmd

import (
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

// subgroupCmd represents the subgroup command
var gitlabEnvsCreateCmd = &cobra.Command{
	Use:   "create {project_id} {env_file_path}",
	Args:  cobra.ExactArgs(2),
	Short: "Create ENVs for Gitlab project",
	Long: `
  Create ENVs for a specific Gitlab project.
  You can also create the env for a specific environment using the 
  flag -e. Please see the example section.
	`,
	Example: `	
  Create ENVs for the project 1234.
  opsi gitlab envs create 1234 /file/to/env.yml

  ---

  Create ENVS for the project 1234 but only for staging environment
  opsi gitlab envs create 1234 /file/to/env.yml -e staging
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
			ui.Fatal(err)
		}
	},
}

func init() {
	gitlabEnvsCmd.AddCommand(gitlabEnvsCreateCmd)
	gitlabEnvsCreateCmd.Flags().StringP("env", "e", "*", "The environment scope")
}
