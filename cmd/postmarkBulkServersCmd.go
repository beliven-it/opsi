package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var postmarkBulkServerCmd = &cobra.Command{
	Use:   "servers",
	Short: "Bulk actions on servers",
	Long:  `Bulk actions on servers`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check for configuration
		// TODO: consider to move inside the scope packet
		if mainConfig.Postmark.SlackWebhook == "" {
			fmt.Println("Missing slack webhook in configuration")
			os.Exit(1)
		}

		postmark.BulkEditServers()
	},
}

func init() {
	postmarkBulkCmd.AddCommand(postmarkBulkServerCmd)
}
