package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gitlabCreateMirrorCmd = &cobra.Command{
	Use:   "mirror {project}...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Enable the mirroring for existing Gitlab projects",
	Long: `Enable the mirroring for one or more existing Gitlab projects.
Each project can be given as numeric ID or as full path.
The destination project is created on the mirror instance when missing.
Projects that already have a mirror are refused: use "update mirroring" for those.`,
	Example: `
  Enable the mirroring for a project
  opsi gitlab create mirror corporate/wiki/beliven-wiki

  ---

  Enable the mirroring for more projects, by path or ID
  opsi gitlab create mirror corporate/wiki/beliven-wiki 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gitlab.CreateMirror(args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateMirrorCmd)
}
