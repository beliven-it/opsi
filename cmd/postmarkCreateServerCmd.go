package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var postmarkCreateServerCmd = &cobra.Command{
	Use:   "server {server_name}",
	Args:  cobra.ExactArgs(1),
	Short: "Create a postmark server",
	Long:  `Create a postmark server`,
	Example: `
  Create a postmark server with name "my-server"
  opsi postmark create my-server	

  ---

  Create a postmark server with name "my-server" and assign the red color
  opsi postmark create my-server -c red
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: condider to move in scope location
		if mainConfig.Postmark.SlackWebhook == "" {
			fmt.Println("Missing slack webhook in configuration")
			os.Exit(1)
		}

		// Take the name of postmark server
		name := args[0]

		// Take the color to attach
		color, _ := cmd.Flags().GetString("color")

		// Create server
		postmark.CreateServer(name, color)
	},
}

func init() {
	postmarkCreateCmd.AddCommand(postmarkCreateServerCmd)
	postmarkCreateServerCmd.Flags().StringP("color", "c", "blue", "The color settings for the server")
}
