/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	slugify "github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var gitlabCreateProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create a Gitlab project",
	Long: `
	Create a Gitlab project. 
	You must provide a valid name as first argument.

	You can also provide a subgroup instead using the default one.

	Make sure to have administrator permission to perform this request.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing name argument")
		}

		subgroup, _ := cmd.Flags().GetInt("subgroup")
		if subgroup == 0 {
			return errors.New("missing subgroup ID")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		subgroup, _ := cmd.Flags().GetInt("subgroup")
		pathname, _ := cmd.Flags().GetString("path")

		if pathname == "" {
			pathname = slugify.Slugify(name)
		}

		err := gitlab.CreateProject(name, pathname, subgroup)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateProjectCmd)

	gitlabCreateProjectCmd.Flags().IntP("subgroup", "s", 0, "The subgroup ID")
	gitlabCreateProjectCmd.Flags().StringP("path", "p", "", "The slugify name for the project")
}
