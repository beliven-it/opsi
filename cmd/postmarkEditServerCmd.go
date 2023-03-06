package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var postmarkEditServerCmd = &cobra.Command{
	Use:     "server {server_id}",
	Args:    cobra.ExactArgs(1),
	Short:   "Edit a postmark server",
	Long:    "Edit a postmark server",
	Example: "TODO....",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Move in different location
		if mainConfig.Postmark.SlackWebhook == "" {
			fmt.Println("Missing slack webhook in configuration")
			os.Exit(1)
		}

		// Take the server ID from arguments and
		// convert into a integer rapresentation.
		idAsNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid ID type. Integer expected")
			os.Exit(1)
		}

		// Edit server
		postmark.EditServer(idAsNumber)
	},
}

func init() {
	postmarkEditCmd.AddCommand(postmarkEditServerCmd)
}
