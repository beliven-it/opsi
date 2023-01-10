package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "The gitlab scope commands",
	Long: `
	You must provide a valid verb.

- deprovisoning to perform the deprovisoning of a user.
- create to perform the creation for the entity specified.
- bulk to perform the execution of a massive action.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing verb argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(gitlabCmd)
}
