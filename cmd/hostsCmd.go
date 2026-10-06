package cmd

import (
	"github.com/spf13/cobra"
)

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Args:  cobra.NoArgs,
	Short: "Check the hosts of the infrastructure",
	Long: `Check the hosts of the infrastructure.
The list of hosts is the one of the hssh CLI.`,
	Run: showHelp,
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
