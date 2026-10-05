package cmd

import (
	"github.com/spf13/cobra"
)

var gitlabSubgroupCmd = &cobra.Command{
	Use:   "subgroup {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "Manage Gitlab subgroups",
	Long:  "Manage Gitlab subgroups",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabSubgroupCmd)
}
