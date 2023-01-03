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
	Long: `
	Update gitlab settings projects.
	Make sure to have administrator permission to perform this request.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := gitlab.BulkSettings()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabBulkCmd.AddCommand(gitlabBulkSettingsCmd)
}
