package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var postmarkBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Allow to perform a bulk action on entity",
	Long: `You can perform the action in one of the following entities:

- servers
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
	postmarkCmd.AddCommand(postmarkBulkCmd)
}
