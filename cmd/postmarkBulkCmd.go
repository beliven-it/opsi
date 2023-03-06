package cmd

import "github.com/spf13/cobra"

var postmarkBulkCmd = &cobra.Command{
	Use:   "bulk {entity}",
	Args:  cobra.ExactArgs(1),
	Short: "Allow to perform a bulk action on entity",
	Long:  "Allow to perform a bulk action on entity",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	postmarkCmd.AddCommand(postmarkBulkCmd)
}
