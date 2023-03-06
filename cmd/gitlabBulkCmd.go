package cmd

import (
	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabBulkCmd = &cobra.Command{
	Use:   "bulk {entity}",
	Args:  cobra.ExactArgs(1),
	Short: "Perform a massive action on an entity",
	Long:  "Perform a massive action on an entity",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabBulkCmd)
}
