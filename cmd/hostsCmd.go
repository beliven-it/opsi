package cmd

import (
	"github.com/spf13/cobra"
)

var hostsCmd = &cobra.Command{
	Use:   "hosts {verb}",
	Args:  cobra.ExactArgs(1),
	Short: "The hosts scope commands",
	Long:  "The hosts scope commands",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
