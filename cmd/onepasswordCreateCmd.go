package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var onepasswordCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Allow to create the necessary assets for project for 1password",
	Long: `In the specific this command generate:

- Vault PRI and relative group to allow to store private assets.
- Vault PUB and relative groups to allow to store public assets.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing project name argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		content, err := Scripts.ReadFile("scripts/1password/create.bash")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, err = onepassword.Create(content, args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	onepasswordCmd.AddCommand(onepasswordCreateCmd)
}
