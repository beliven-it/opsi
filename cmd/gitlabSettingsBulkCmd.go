package cmd

import (
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var gitlabSettingsBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Apply the default settings to all Gitlab projects",
	Long:  "Apply the default settings to all Gitlab projects.",
	Example: `
  Apply the default settings to all projects
  opsi gitlab settings bulk
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create the output channel for the messages
		channel := make(chan string)
		printed := make(chan struct{})
		go func() {
			defer close(printed)
			for item := range channel {
				ui.Print("%s", item)
			}
		}()

		// Execute bulk
		err := gitlab.BulkSettings(&channel)
		close(channel)
		<-printed
		if err != nil {
			ui.Fatal(err)
		}
	},
}

func init() {
	gitlabSettingsCmd.AddCommand(gitlabSettingsBulkCmd)
}
