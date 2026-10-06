package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabMirrorCmd = &cobra.Command{
	Use:   "mirror",
	Args:  cobra.NoArgs,
	Short: "Manage the mirroring of Gitlab projects",
	Long:  "Manage the mirroring of Gitlab projects",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabMirrorCmd)
}
