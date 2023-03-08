package cmd

import (
	"fmt"
	"opsi/helpers"
	"os"

	"github.com/spf13/cobra"
)

var gitlabDeprovisioningCmd = &cobra.Command{
	Use:   "deprovisioning {username}",
	Args:  cobra.ExactArgs(1),
	Short: "Remove an user from all groups and projects",
	Long:  "Remove an user from all groups and projects",
	Example: `
  Remove the user john.doe from gitlab.
  opsi gitlab deprovisioning john.doe	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take the username
		username := args[0]

		// Confirm the action
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			helpers.Confirm()
		}

		// Deprovisioning the user
		err := gitlab.Deprovionioning(username)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabCmd.AddCommand(gitlabDeprovisioningCmd)
	gitlabDeprovisioningCmd.Flags().BoolP("force", "f", false, "Not ask confirmation to delete")
}
