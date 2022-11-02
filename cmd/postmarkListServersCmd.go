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
var postmarkListServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "Return a list of postmark servers",
	Long:  `Return a list of postmark servers`,
	Run: func(cmd *cobra.Command, args []string) {
		servers, err := postmark.GetServers()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for index, server := range servers {
			fmt.Printf("[%d] %s (%d)\n", index, server.Name, server.ID)
		}
	},
}

func init() {
	postmarkCmd.AddCommand(postmarkListServersCmd)
}
