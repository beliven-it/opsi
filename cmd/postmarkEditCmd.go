package cmd

import (
	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var postmarkEditCmd = &cobra.Command{
	Use:   "edit {entity}",
	Args:  cobra.ExactArgs(1),
	Short: "Edit a postmark entity",
	Long:  "Edit a postmark entity",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	postmarkCmd.AddCommand(postmarkEditCmd)
}
