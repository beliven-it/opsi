/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var postmarkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Allow to create a specific entity",
	Long: `You can create one of the following entities:

	- server
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
	postmarkCmd.AddCommand(postmarkCreateCmd)
}
