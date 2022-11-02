/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var postmarkListCmd = &cobra.Command{
	Use:   "list",
	Short: "Allow to list entities",
	Long: `You can list one of the following entities:

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
	postmarkCmd.AddCommand(postmarkListCmd)
}
