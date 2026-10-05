package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabMirrorCmd = &cobra.Command{
	Use:   "mirror {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage the mirroring of Gitlab projects",
	Long:  "Manage the mirroring of Gitlab projects",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabMirrorCmd)
}
