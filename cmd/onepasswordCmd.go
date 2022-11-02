/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// postmarkCmd represents the postmark command
var onepasswordCmd = &cobra.Command{
	Use:   "1password",
	Short: "The 1password scope commands",
	Long: `
	You must provide a valid verb.

- create to perform the creation for the entity specified.
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
	rootCmd.AddCommand(onepasswordCmd)
}
