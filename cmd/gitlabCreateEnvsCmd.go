package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// subgroupCmd represents the subgroup command
var gitlabCreateEnvsCmd = &cobra.Command{
	Use:   "envs",
	Short: "Create ENVs for Gitlab project",
	Long: `Create ENVs for Gitlab project. 
	You must provide a valid project ID.
	Make sure to have administrator permission to perform this request.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing Project ID argument")
		}

		if len(args) == 1 {
			return errors.New("missing Env file argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		envFile := args[1]

		env, _ := cmd.Flags().GetString("env")

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
