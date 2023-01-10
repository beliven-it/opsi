package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Allow to delete a specific entity",
	Long: `You can delete one of the following entities:

- envs
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing entity argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	gitlabCmd.AddCommand(gitlabDeleteCmd)
}
