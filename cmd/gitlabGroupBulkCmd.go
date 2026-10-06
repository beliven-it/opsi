package cmd

import (
	"fmt"
	"opsi/helpers"
	"opsi/helpers/ui"
	"sort"

	"github.com/spf13/cobra"
)

var gitlabGroupBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Align gitlab groups to the default settings",
	Long: `Align gitlab groups to the default settings.

The command always reports the groups that differ before touching anything,
and asks for confirmation before applying.`,
	Example: `
  Report the groups that differ, without applying anything
  opsi gitlab group bulk --dry-run

  ---

  Align every group, asking for confirmation first
  opsi gitlab group bulk

  ---

  Align a single group and its subgroups
  opsi gitlab group bulk -g 1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		group, _ := cmd.Flags().GetInt("group")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")

		drifts, err := gitlab.GroupSettingsDrift(group)
		if err != nil {
			ui.Fatal(err)
		}

		if len(drifts) == 0 {
			ui.Success("Every group already matches the default settings")
			return
		}

		for _, drift := range drifts {
			ui.Section("%s %s", drift.FullPath, ui.Dim(fmt.Sprintf("#%d", drift.ID)))

			keys := make([]string, 0, len(drift.Changes))
			for key := range drift.Changes {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				ui.Item("%s: %v -> %v", key, drift.From[key], drift.Changes[key])
			}
		}

		label := "groups"
		if len(drifts) == 1 {
			label = "group"
		}
		ui.Section("%d %s to update", len(drifts), label)

		if dryRun {
			return
		}

		if !force {
			helpers.Confirm()
		}

		channel := make(chan string)
		printed := make(chan struct{})
		go func() {
			defer close(printed)
			for item := range channel {
				ui.Print("%s", item)
			}
		}()

		err = gitlab.ApplyGroupSettings(drifts, &channel)
		close(channel)
		<-printed
		if err != nil {
			ui.Fatal(err)
		}
	},
}

func init() {
	gitlabGroupCmd.AddCommand(gitlabGroupBulkCmd)
	gitlabGroupBulkCmd.Flags().IntP("group", "g", 0, "Restrict the update to a group and its subgroups")
	gitlabGroupBulkCmd.Flags().Bool("dry-run", false, "Only report the groups that differ")
	gitlabGroupBulkCmd.Flags().BoolP("force", "f", false, "Not ask confirmation before applying")
}
