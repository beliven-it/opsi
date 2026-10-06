package cmd

import (
	"fmt"
	"opsi/helpers"
	"os"

	"github.com/spf13/cobra"
)

var gitlabUserDeprovisioningCmd = &cobra.Command{
	Use:   "deprovisioning {username|email}",
	Args:  cobra.ExactArgs(1),
	Short: "Remove a user from all groups and projects",
	Long:  "Remove a user from all groups and projects",
	Example: `
  Remove the user john.doe from gitlab.
  opsi gitlab user deprovisioning john.doe

  ---

  The user can be searched by email too.
  opsi gitlab user deprovisioning john.doe@example.com

  ---

  Show where john.doe is member, without removing anything.
  opsi gitlab user deprovisioning john.doe --dry-run
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take the username
		username := args[0]

		// Confirm the action (a dry run changes nothing)
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		if !force && !dryRun {
			helpers.Confirm()
		}

		// Deprovisioning the user
		err := gitlab.Deprovisioning(username, dryRun)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabUserCmd.AddCommand(gitlabUserDeprovisioningCmd)
	gitlabUserDeprovisioningCmd.Flags().Bool("dry-run", false, "Only list the memberships that would be removed")
	gitlabUserDeprovisioningCmd.Flags().BoolP("force", "f", false, "Not ask confirmation to delete")
}
