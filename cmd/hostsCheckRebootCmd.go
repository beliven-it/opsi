package cmd

import (
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

var hostsCheckRebootCmd = &cobra.Command{
	Use:   "check-reboot",
	Short: "Check which hosts need to reboot",
	Long:  "Check which hosts need to reboot. The list of hosts is the one of the hssh CLI.",
	Example: `
  Check all the hosts
  opsi hosts check-reboot
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check reboot
		err := hosts.CheckReboot()
		if err != nil {
			ui.Fatal(err)
		}
	},
}

func init() {
	hostsCmd.AddCommand(hostsCheckRebootCmd)
}
