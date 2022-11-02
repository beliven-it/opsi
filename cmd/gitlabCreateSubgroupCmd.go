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

// subgroupCmd represents the subgroup command
var gitlabCreateSubgroupCmd = &cobra.Command{
	Use:   "subgroup",
	Short: "Create a Gitlab subgroup",
	Long: `Create a Gitlab subgroup. 
	You must provide a valid name.
	Make sure to have administrator permission to perform this request.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing name argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		parent, _ := cmd.Flags().GetInt("parent")
		pathname, _ := cmd.Flags().GetString("path")

		if parent == 0 {
			parent = mainConfig.Gitlab.GroupID
		}

		if pathname == "" {
			pathname = slugify.Slugify(name)
		}

		err := gitlab.CreateSubgroup(name, pathname, parent)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	gitlabCreateCmd.AddCommand(gitlabCreateSubgroupCmd)

	gitlabCreateSubgroupCmd.Flags().IntP("parent", "s", 0, "The parent of the subgroup")
	gitlabCreateSubgroupCmd.Flags().StringP("path", "p", "", "The slugify name for the subgroup")
}
