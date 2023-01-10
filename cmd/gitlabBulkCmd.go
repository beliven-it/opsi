package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Perform a massive action on an entity",
	Long: `You can perform this action on one of the following entities:

	- settings
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing entity argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitlab bulk")
	},
}

func init() {
	gitlabCmd.AddCommand(gitlabBulkCmd)
}
