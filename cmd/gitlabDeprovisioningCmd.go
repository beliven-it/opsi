/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deprovisioningCmd represents the deprovisioning command
var gitlabDeprovisioningCmd = &cobra.Command{
	Use:   "deprovisioning",
	Short: "Deprovisioning an user",
	Long: `Deprovisioning an user. 
	You must provide a valid username.
	Make sure to have administrator permission to perform this request.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing username argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]

		err := gitlab.Deprovionioning(username)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	gitlabCmd.AddCommand(gitlabDeprovisioningCmd)
}
