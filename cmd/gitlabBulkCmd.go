/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var gitlabBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitlab called")
	},
}

func init() {
	gitlabCmd.AddCommand(gitlabBulkCmd)
}
