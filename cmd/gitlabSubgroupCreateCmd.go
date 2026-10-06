package cmd

import (
	"opsi/helpers/ui"

	slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

var gitlabSubgroupCreateCmd = &cobra.Command{
	Use:   "create {subgroup_name}",
	Args:  cobra.ExactArgs(1),
	Short: "Create a Gitlab subgroup",
	Long:  "Create a Gitlab subgroup",
	Example: `
  Create a subgroup with name "research" attach to a specific group with id 1234
  opsi gitlab subgroup create research -s 1234 

  ---

  Create a subgroup with name "development" but with path to "devs"
  opsi gitlab subgroup create development -p devs -s 1234
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
			ui.Fatal(err)
		}

		ui.Success("Created new subgroup with ID %d", subgroupID)
	},
}

func init() {
	gitlabSubgroupCmd.AddCommand(gitlabSubgroupCreateCmd)
	gitlabSubgroupCreateCmd.Flags().IntP("parent", "s", 0, "The parent of the subgroup you want create")
	gitlabSubgroupCreateCmd.Flags().StringP("path", "p", "", "The slugify name for the subgroup")

	// Mark group as required
	gitlabSubgroupCreateCmd.MarkFlagRequired("parent")
}
