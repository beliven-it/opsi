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
		err := hosts.CheckReboot()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	},
}

func init() {
	hostsCmd.AddCommand(hostsCheckRebootCmd)
}
