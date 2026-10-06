package cmd

import (
	"opsi/helpers/ui"

	gl "opsi/scopes/gitlab"

	slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var gitlabProjectCreateCmd = &cobra.Command{
	Use:   "create {project_name}",
	Args:  cobra.ExactArgs(1),
	Short: "Create a Gitlab project",
	Long:  "This command allow to create a Gitlab project in a specific workspace",
	Example: `
  Create a project with name "Password manager" for subgroup 12345:
  opsi gitlab project create "Password manager" -s 12345 

  ---

  Create a project with name "Delorian" but path "my-delorian"
  opsi gitlab project create Delorian -p my-delorian

  ---

  Create a project with name "Valerian" using "master" as default branch
  opsi gitlab project create Valerian -b master

  ---

  Create a project with name "Akkadian" enabling mirroring to another gitlab
  opsi gitlab project create Akkadian -m

  ---

  Create a project with name "Nightly" with the default weekly pipeline schedule
  opsi gitlab project create Nightly -c

  ---

  Create a project with name "Anonymous" disabling shared runners
  opsi gitlab project create Anonymous -r

  ---

  Create a project with visibility "internal"
  opsi gitlab project create Anonymous -i internal

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
		schedule, _ := cmd.Flags().GetBool("schedule")

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
			Schedule:      schedule,
			Group:         group,
		}

		// Create the project
		projectID, err := gitlab.CreateProject(payload)

		if err != nil {
			ui.Fatal(err)
		}

		ui.Success("Created new project with ID %d", projectID)
	},
}

func init() {
	gitlabProjectCmd.AddCommand(gitlabProjectCreateCmd)
	gitlabProjectCreateCmd.Flags().IntP("group", "s", 0, "the group associated to the project. If not provided the one in the configuration will be used")
	gitlabProjectCreateCmd.Flags().StringP("path", "p", "", "the path for the project. This flag is useful if you don't want to use the project name for the path")
	gitlabProjectCreateCmd.Flags().StringP("branch-default", "b", "main", "the default main branch. Possible values are master or main")
	gitlabProjectCreateCmd.Flags().BoolP("mirror", "m", false, "Enable or disable the mirroring repo. Default is false")
	gitlabProjectCreateCmd.Flags().BoolP("sharedrunners", "r", false, "Enable or disable the shared runners. Default is true")
	gitlabProjectCreateCmd.Flags().BoolP("schedule", "c", false, "Create the default weekly pipeline schedule, on the least used slot. Default is false")
	gitlabProjectCreateCmd.Flags().StringP("visibility", "i", "", "Set the visibility of the project. Allowed values are private, public, internal")

	// Mark group as required
	gitlabProjectCreateCmd.MarkFlagRequired("group")
}
