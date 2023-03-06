package cmd

import (
	"github.com/spf13/cobra"
)

var postmarkCmd = &cobra.Command{
	Use:   "postmark {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "The postmark commands",
	Long:  "The postmark commands",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(postmarkCmd)
}
