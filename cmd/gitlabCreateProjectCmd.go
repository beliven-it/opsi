package cmd

import (
	"fmt"
	"os"

	gl "opsi/scopes/gitlab"

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
  opsi gitlab create project "Password manager" -s 12345 

  ---

  Create a project with name "Delorian" but path "my-delorian"
  opsi gitlab create project Delorian -p my-delorian

  ---

  Create a project with name "Valerian" using "master" as default branch
  opsi gitlab create project Valerian -b master

  ---

  Create a project with name "Akkadian" enabling mirroring to another gitlab
  opsi gitlab create project Akkadian -m

  ---

  Create a project with name "Anonymous" disabling shared runners
  opsi gitlab create project Anonymous -r

  ---

  Create a project with visibility "internal"
  opsi gitlab create project Anonymous -i internal

	`,

	Run: func(cmd *cobra.Command, args []string) {
		// Take project name
		name := args[0]

		// Take flags
		group, _ := cmd.Flags().GetInt("group")
		pathname, _ := cmd.Flags().GetString("path")
		defaultBranch, _ := cmd.Flags().GetString("branch-default")
		mirror, _ := cmd.Flags().GetBool("mirror")
		sharedRunners, _ := cmd.Flags().GetBool("sharedrunners")
		visibility, _ := cmd.Flags().GetString("visibility")

		// Slugify the name if the pathname flag
		// for the project is not provided
		if pathname == "" {
			pathname = slugify.Slugify(name)
		}

		// Prepare the payload
		payload := gl.ProjectRequest{
			Name:          name,
			Path:          pathname,
			Visibility:    visibility,
			DefaultBranch: defaultBranch,
			Mirror:        mirror,
			SharedRunners: sharedRunners,
			Group:         group,
		}

		// Create the project
		projectID, err := gitlab.CreateProject(payload)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Created new project with ID", projectID)
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateProjectCmd)
	gitlabCreateProjectCmd.Flags().IntP("group", "s", 0, "the group associated to the project. If not provided the one in the configuration will be used")
	gitlabCreateProjectCmd.Flags().StringP("path", "p", "", "the path for the project. This flag is useful if you don't want to use the project name for the path")
	gitlabCreateProjectCmd.Flags().StringP("branch-default", "b", "main", "the default main branch. Possible values are master or main")
	gitlabCreateProjectCmd.Flags().BoolP("mirror", "m", false, "Enable or disable the mirroring repo. Default is false")
	gitlabCreateProjectCmd.Flags().BoolP("sharedrunners", "r", false, "Enable or disable the shared runners. Default is true")
	gitlabCreateProjectCmd.Flags().StringP("visibility", "i", "", "Set the visibility of the project. Allowed values are private, public, internal")

	// Mark group as required
	gitlabCreateProjectCmd.MarkFlagRequired("group")
}
