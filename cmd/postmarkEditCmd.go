package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var postmarkEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a postmark server",
	Long:  `Edit a postmark server`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing entity argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	postmarkCmd.AddCommand(postmarkEditCmd)
}
