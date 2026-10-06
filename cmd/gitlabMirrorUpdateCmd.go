package cmd

import (
	"opsi/helpers/ui"

	// gl "opsi/scopes/gitlab"

	// slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

// updateMirroring represents the update mirroring command
var gitlabMirrorUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the mirroring of all Gitlab projects",
	Long:  "Update the mirroring of all Gitlab projects that have one.",
	Example: `
  Update the mirroring of all projects
  opsi gitlab mirror update
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Update mirroring
		err := gitlab.UpdateMirroring()
		if err != nil {
			ui.Fatal(err)
		}

		ui.Success("All GitLab repositories have been successfully updated")
	},
}

func init() {
	gitlabMirrorCmd.AddCommand(gitlabMirrorUpdateCmd)
}
