package cmd

import (
	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "The gitlab scope commands",
	Long:  "The gitlab scope commands",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(gitlabCmd)
}
