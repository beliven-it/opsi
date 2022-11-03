/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// postmarkCmd represents the postmark command
var postmarkCmd = &cobra.Command{
	Use:   "postmark",
	Short: "The postmark commands",
	Long: `You must provide a valid verb.

	- list to list a specified entity.
	- create to perform the creation for the entity specified.
	- edit to perform the execution for update the entity specified`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("postmark called")
	},
}

func init() {
	rootCmd.AddCommand(postmarkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postmarkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postmarkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
