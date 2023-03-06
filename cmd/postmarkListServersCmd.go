package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var postmarkListServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "Return a list of postmark servers",
	Long:  `Return a list of postmark servers`,
	Example: `
  Take the list of the servers:
  opsi postmark list servers	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take the servers
		servers, err := postmark.GetServers()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// Render the list
		for index, server := range servers {
			fmt.Printf("[%d] %s (%d)\n", index, server.Name, server.ID)
		}
	},
}

func init() {
	postmarkListCmd.AddCommand(postmarkListServersCmd)
}
