package cmd

import (
	"fmt"
	"os"

	// gl "opsi/scopes/gitlab"

	// slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

// updateMirroring represents the update mirroring command
var gitlabUpdateMirroringCmd = &cobra.Command{
	Use:   "mirroring",
	Short: "Update Gitlab Mirroring",
	Long:  "This command updates mirroring for all GitLab repositories",

	Run: func(cmd *cobra.Command, args []string) {
		// Update mirroring
		err := gitlab.UpdateMirroring()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("All GitLab repositories have been successfully updated")
	},
}

func init() {
	gitlabUpdateCmd.AddCommand(gitlabUpdateMirroringCmd)
}
