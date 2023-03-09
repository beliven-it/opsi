package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var onepasswordCreateCmd = &cobra.Command{
	Use:   "create {project_name}",
	Args:  cobra.ExactArgs(1),
	Short: "Allow to create a project inside a 1password environment",
	Long: `
Allow to create a project inside a 1password environment.
In specific, the following entities will be created:

- Vault PRI and relative group to allow to store private assets.
- Vault PUB and relative groups to allow to store public assets.
	`,
	Example: `  
  Create a 1password vault called "personal vault"
  opsi 1password create "personal vault"	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := onepassword.Create(args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	onepasswordCmd.AddCommand(onepasswordCreateCmd)
}
