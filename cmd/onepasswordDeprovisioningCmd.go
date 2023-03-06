package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var onepasswordDeprovisioningCmd = &cobra.Command{
	Use:   "deprovisioning",
	Short: "Deprovision 1password inactive users",
	Long: `Deprovision 1password inactive users. 
	If you need to deprovisioning a specific user, you can use the -e flag
	and search the user by email. 

	Show the examples and flags sections for further informations
	`,
	Example: `
  Deprovisioning all inactive users from 1password workspace
  opsi 1password deprovisioning	

  ---

  Deprovisioning the user with email john.doe@example.com from 1password workspace
  opsi 1password deprovisioning -e john.doe@example.com
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Take email from flag
		email, _ := cmd.Flags().GetString("email")

		// TODO: add force

		// Start deprovisioning
		err := onepassword.Deprovisioning(email)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	onepasswordCmd.AddCommand(onepasswordDeprovisioningCmd)
	onepasswordDeprovisioningCmd.Flags().StringP("email", "e", "", "The email of the user to deprovisioning")
}
