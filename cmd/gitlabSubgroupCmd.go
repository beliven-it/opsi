package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabSubgroupCmd = &cobra.Command{
	Use:   "subgroup",
	Args:  cobra.NoArgs,
	Short: "Manage Gitlab subgroups",
	Long:  "Manage Gitlab subgroups",
	Run:   showHelp,
}

func init() {
	gitlabCmd.AddCommand(gitlabSubgroupCmd)
}
