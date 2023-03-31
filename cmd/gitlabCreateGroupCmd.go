package cmd

import (
	"fmt"
	"os"

	slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

var gitlabCreateGroupCmd = &cobra.Command{
	Use:   "group {group_name}",
	Args:  cobra.ExactArgs(1),
	Short: "Create a Gitlab group",
	Long:  "Create a Gitlab group",
	Example: `
  Create a group with name "research"
  opsi gitlab create group research

  ---

  Create a group with name "development" but with path to "devs"
  opsi gitlab create group development -p devs

  ---

  Create a group with name "development" with "public" visibility
  opsi gitlab create group development -i public
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take the name of the group
		name := args[0]

		// Take the pathname from the flag
		pathname, _ := cmd.Flags().GetString("path")

		// Take the visibility from the flag
		visibility, _ := cmd.Flags().GetString("visibility")

		// If the pathname is not provided
		// let the system slugify the name
		if pathname == "" {
			pathname = slugify.Slugify(name)
		}

		// Create subgroup
		groupID, err := gitlab.CreateGroup(name, pathname, visibility)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Created new group with ID", groupID)
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateGroupCmd)
	gitlabCreateGroupCmd.Flags().StringP("path", "p", "", "The slugify name for the group")
	gitlabCreateGroupCmd.Flags().StringP("visibility", "i", "private", "Set the visibility of the group")
}
