package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var onepasswordDeprovisioningCmd = &cobra.Command{
	Use:   "deprovisioning",
	Short: "Deprovision 1password users",
	Long:  "Deprovision 1password users. You can also specify an user email to deprovisioning only selected user.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check and assign email argument
		email := ""
		if len(args) > 0 {
			email = args[0]
		}

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
}
