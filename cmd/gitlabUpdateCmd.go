package cmd

import (
	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabUpdateCmd = &cobra.Command{
	Use:   "update {entity}",
	Args:  cobra.ExactArgs(1),
	Short: "Allow to update a specific entity",
	Long:  "Allow to update a specific entity",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabUpdateCmd)
}
