package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabUserCmd = &cobra.Command{
	Use:   "user",
	Args:  cobra.NoArgs,
	Short: "Manage Gitlab users",
	Long:  "Manage Gitlab users",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabUserCmd)
}
