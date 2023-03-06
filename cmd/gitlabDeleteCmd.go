package cmd

import (
	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabDeleteCmd = &cobra.Command{
	Use:   "delete {entity}",
	Args:  cobra.ExactArgs(1),
	Short: "Allow to delete a specific entity",
	Long:  "Allow to delete a specific entity",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabDeleteCmd)
}
