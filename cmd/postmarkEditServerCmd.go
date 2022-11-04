/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var postmarkEditServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Edit a postmark server",
	Long:  `Edit a postmark server`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing ID argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if mainConfig.Postmark.SlackWebhook == "" {
			fmt.Println("Missing slack webhook in configuration")
			os.Exit(1)
		}

		idAsNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid ID type. Integer expected")
			os.Exit(1)
		}

		postmark.EditServer(idAsNumber)
	},
}

func init() {
	postmarkEditCmd.AddCommand(postmarkEditServerCmd)
}
