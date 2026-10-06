package cmd

import (
	"fmt"
	"opsi/helpers"
	"os"

	"github.com/spf13/cobra"
)

var onepasswordUserDeprovisioningCmd = &cobra.Command{
	Use:   "deprovisioning",
	Short: "Deprovision 1password inactive users",
	Long: `Deprovision 1password inactive users. 
	If you need to deprovisioning a specific user, you can use the -e flag
	and search the user by email. 

	Show the examples and flags sections for further informations
	`,
	Example: `
  Deprovisioning all inactive users from 1password workspace
  opsi 1password user deprovisioning	

  ---

  Deprovisioning the user with email john.doe@example.com from 1password workspace
  opsi 1password user deprovisioning -e john.doe@example.com

  ---

  Show which suspended users would be deleted, without deleting them
  opsi 1password user deprovisioning --dry-run
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take email from flag
		email, _ := cmd.Flags().GetString("email")

		// Confirm the action (a dry run changes nothing)
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		if !force && !dryRun {
			helpers.Confirm()
		}

		// Start the deprovisioning procedure
		err := onepassword.Deprovisioning(email, dryRun)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	onepasswordUserCmd.AddCommand(onepasswordUserDeprovisioningCmd)
	onepasswordUserDeprovisioningCmd.Flags().StringP("email", "e", "", "The email of the user to deprovisioning")
	onepasswordUserDeprovisioningCmd.Flags().Bool("dry-run", false, "Only list the users that would be deleted")
	onepasswordUserDeprovisioningCmd.Flags().BoolP("force", "f", false, "Not ask confirmation to delete")
}
