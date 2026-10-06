package cmd

import (
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var gitlabSettingsBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Update gitlab settings projects",
	Long:  "Update gitlab settings projects",
	Example: `
  Update all projects
  opsi gitlab settings bulk

  Update all projects
  opsi gitlab settings bulk
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create the output channel for the messages
		channel := make(chan string)
		go func() {
			for item := range channel {
				ui.Print("%s", item)
			}
		}()

		// Execute bulk
		err := gitlab.BulkSettings(&channel)
		if err != nil {
			ui.Fatal(err)
		}
	},
}

func init() {
	gitlabSettingsCmd.AddCommand(gitlabSettingsBulkCmd)
}
