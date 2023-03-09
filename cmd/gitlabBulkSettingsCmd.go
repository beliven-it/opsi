package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var gitlabBulkSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Update gitlab settings projects",
	Long:  "Update gitlab settings projects",
	Example: `
  Update all projects
  opsi gitlab bulk settings

  Update all projects
  opsi gitlab bulk settings
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create the output channel for the messages
		channel := make(chan string)
		go func() {
			for item := range channel {
				fmt.Println(item)
			}
		}()

		// Execute bulk
		err := gitlab.BulkSettings(&channel)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabBulkCmd.AddCommand(gitlabBulkSettingsCmd)
}
