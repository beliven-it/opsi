package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabListEnvsCmd = &cobra.Command{
	Use:   "envs {project_id}",
	Args:  cobra.ExactArgs(1),
	Short: "List ENVs for Gitlab project",
	Long:  "List ENVs for Gitlab project",
	Example: `
  Show all envs for the project 1234
  opsi gitlab list envs 1234
  
  ---

  Show all envs for the project 1234 but only for staging environment
  opsi gitlab list env 1234 -e staging
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take project ID
		projectID := args[0]

		// Take env from flag
		env, _ := cmd.Flags().GetString("env")

		// List the envs
		err := gitlab.ListEnvs(projectID, env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabListCmd.AddCommand(gitlabListEnvsCmd)
	gitlabListEnvsCmd.Flags().StringP("env", "e", "*", "The environment scope")
}
