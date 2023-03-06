package cmd

import (
	"fmt"
	"os"

	slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var gitlabCreateProjectCmd = &cobra.Command{
	Use:   "project {project_name}",
	Args:  cobra.ExactArgs(1),
	Short: "Create a Gitlab project",
	Long:  "This command allow to create a Gitlab project in a specific workspace",
	Example: `
  Create a project with name "Password manager" for subgroup 12345:
  opsi gitlab create project "Password manager" -g 12345 

  ---

  Create a project with name "Delorian" but path "my-delorian"
  opsi gitlab create project Delorian -p my-delorian
	`,

	Run: func(cmd *cobra.Command, args []string) {
		// Take project name
		name := args[0]

		// Take flags
		group, _ := cmd.Flags().GetInt("group")
		pathname, _ := cmd.Flags().GetString("path")

		// Slugify the name if the pathname flag
		// for the project is not provided
		if pathname == "" {
			pathname = slugify.Slugify(name)
		}

		// Use the global group ID for the project
		// if the group id flag is not provided
		if group == 0 {
			group = mainConfig.Gitlab.GroupID
		}

		// Create the project
		err := gitlab.CreateProject(name, pathname, group)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateProjectCmd)
	gitlabCreateProjectCmd.Flags().IntP("group", "g", 0, "the group associated to the project. If not provided the one in the configuration will be used")
	gitlabCreateProjectCmd.Flags().StringP("path", "p", "", "the path for the project. This flag is useful if you don't want to use the project name for the path")
}
