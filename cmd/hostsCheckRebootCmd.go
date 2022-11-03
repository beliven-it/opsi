/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the gitlab command
var hostsCheckRebootCmd = &cobra.Command{
	Use:   "check-reboot",
	Short: "Check hosts need to reboot",
	Long: `This command check all the hosts contained in hssh CLI
and check if the machine need to reboot
	`,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := Scripts.ReadFile("scripts/hosts/check-reboot.bash")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		output, err := hosts.CheckReboot(content)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(string(output))
	},
}

func init() {
	hostsCmd.AddCommand(hostsCheckRebootCmd)
}
