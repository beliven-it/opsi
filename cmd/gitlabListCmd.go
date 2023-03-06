package cmd

import (
	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabListCmd = &cobra.Command{
	Use:   "list {entity}",
	Args:  cobra.ExactArgs(1),
	Short: "Show a list of elements for a specific entity",
	Long:  "Show a list of elements for a specific entity",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabListCmd)
}
