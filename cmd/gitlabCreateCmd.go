package cmd

import (
	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabCreateCmd = &cobra.Command{
	Use:   "create {entity}",
	Args:  cobra.ExactArgs(1),
	Short: "Allow to create a specific entity",
	Long:  "Allow to create a specific entity",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabCreateCmd)
}
