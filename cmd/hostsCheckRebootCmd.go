package cmd

import (
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

var hostsCheckRebootCmd = &cobra.Command{
	Use:   "check-reboot",
	Short: "Check hosts need to reboot",
	Long:  "Check hosts need to reboot. The list of hosts are the ones of hssh CLI",
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
