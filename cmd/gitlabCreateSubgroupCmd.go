package cmd

import (
	"fmt"
	"os"

	slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

var gitlabCreateSubgroupCmd = &cobra.Command{
	Use:   "subgroup {subgroup_name}",
	Args:  cobra.ExactArgs(1),
	Short: "Create a Gitlab subgroup",
	Long:  "Create a Gitlab subgroup",
	Example: `
  Create a subgroup with name "research" attach to a specific group with id 1234
  opsi gitlab create subgroup research -s 1234 

  ---

  Create a subgroup with name "development" but with path to "devs"
  opsi gitlab create subgroup development -p devs -s 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take the name of the group
		name := args[0]

		// Take the parent from the flag
		parent, _ := cmd.Flags().GetInt("parent")

		// Take the pathname from the flag
		pathname, _ := cmd.Flags().GetString("path")

		// If the pathname is not provided
		// let the system slugify the name
		if pathname == "" {
			pathname = slugify.Slugify(name)
		}

		// Set the parent as nil if not provided
		var parentAsPointer *int = &parent
		if parent == 0 {
			parentAsPointer = nil
		}

		// Create subgroup
		subgroupID, err := gitlab.CreateSubgroup(name, pathname, parentAsPointer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Created new subgroup with ID", subgroupID)
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateSubgroupCmd)
	gitlabCreateSubgroupCmd.Flags().IntP("parent", "s", 0, "The parent of the subgroup you want create")
	gitlabCreateSubgroupCmd.Flags().StringP("path", "p", "", "The slugify name for the subgroup")

	// Mark group as required
	gitlabCreateSubgroupCmd.MarkFlagRequired("parent")
}
