/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var postmarkEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a postmark server",
	Long:  `Edit a postmark server`,
	Run: func(cmd *cobra.Command, args []string) {
		if mainConfig.Postmark.SlackWebhook == "" {
			fmt.Println("Missing slackwebhook in configuration")
			os.Exit(1)
		}

		postmark.GetServers()
	},
}

func init() {
	postmarkCmd.AddCommand(postmarkEditCmd)
}
