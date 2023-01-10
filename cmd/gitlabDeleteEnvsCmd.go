package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// subgroupCmd represents the subgroup command
var gitlabDeleteEnvsCmd = &cobra.Command{
	Use:   "envs",
	Short: "Delete ENVs for Gitlab project",
	Long: `Delete ENVs for Gitlab project. 
	You must provide a valid project ID.
	Make sure to have administrator permission to perform this request.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing Project ID argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]

		env, _ := cmd.Flags().GetString("env")
		force, _ := cmd.Flags().GetBool("force")

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
