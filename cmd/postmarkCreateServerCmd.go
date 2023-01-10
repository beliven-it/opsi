package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var postmarkCreateServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Create a postmark server",
	Long:  `Create a postmark server`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing name argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if mainConfig.Postmark.SlackWebhook == "" {
			fmt.Println("Missing slack webhook in configuration")
			os.Exit(1)
		}

		name := args[0]
		color, _ := cmd.Flags().GetString("color")

		postmark.CreateServer(name, color)
	},
}

func init() {
	postmarkCreateCmd.AddCommand(postmarkCreateServerCmd)

	postmarkCreateServerCmd.Flags().StringP("color", "c", "blue", "The color settings for the server")
}
