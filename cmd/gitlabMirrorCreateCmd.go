package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabMirrorCreateCmd = &cobra.Command{
	Use:   "create {project}...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Enable the mirroring for existing Gitlab projects",
	Long: `Enable the mirroring for one or more existing Gitlab projects.
Each project can be given as numeric ID or as full path.
The destination project is created on the mirror instance when missing.
Projects that already have a mirror are refused: use "mirror update" for those.`,
	Example: `
  Enable the mirroring for a project
  opsi gitlab mirror create corporate/wiki/beliven-wiki

  ---

  Enable the mirroring for more projects, by path or ID
  opsi gitlab mirror create corporate/wiki/beliven-wiki 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gitlab.CreateMirror(args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabMirrorCmd.AddCommand(gitlabMirrorCreateCmd)
}
