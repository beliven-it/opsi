/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "The hosts scope commands",
	Long: `
	You must provide a valid verb.

- check-reboot to perform a search of hosts need to reboot.
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
	rootCmd.AddCommand(hostsCmd)
}
