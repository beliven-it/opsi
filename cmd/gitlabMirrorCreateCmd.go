package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabMirrorCreateCmd = &cobra.Command{
	Use:   "create {project_id}...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Enable the mirroring for existing Gitlab projects",
	Long: `Enable the mirroring for one or more existing Gitlab projects.
Each project is given by its numeric ID.
The destination project is created on the mirror instance when missing.
Projects that already have a mirror are refused: use "mirror update" for those.`,
	Example: `
  Enable the mirroring for a project
  opsi gitlab mirror create 239

  ---

  Enable the mirroring for more projects
  opsi gitlab mirror create 239 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := projectIDs(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := gitlab.CreateMirror(ids); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabMirrorCmd.AddCommand(gitlabMirrorCreateCmd)
}
