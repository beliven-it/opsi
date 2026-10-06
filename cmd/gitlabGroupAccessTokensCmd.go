package cmd

import (
	"opsi/helpers/ui"

	"github.com/spf13/cobra"
)

var gitlabGroupAccessTokensCmd = &cobra.Command{
	Use:   "access-tokens",
	Args:  cobra.NoArgs,
	Short: "List the access tokens of all Gitlab groups",
	Long: `List the access tokens of all Gitlab groups, with name, scopes and expiration.
	Use --expiring to show only the expired tokens and the ones expiring soon.
	The token must be able to read the access tokens of every group (admin).`,
	Example: `
  List the access tokens of all groups
  opsi gitlab group access-tokens

  ---

  List the tokens already expired or expiring within 30 days
  opsi gitlab group access-tokens --expiring 30
	`,
	Run: func(cmd *cobra.Command, args []string) {
		expiring, _ := cmd.Flags().GetInt("expiring")

		err := gitlab.ListGroupAccessTokens(expiring)
		if err != nil {
			ui.Fatal(err)
		}
	},
}

func init() {
	gitlabGroupCmd.AddCommand(gitlabGroupAccessTokensCmd)
	gitlabGroupAccessTokensCmd.Flags().Int("expiring", -1, "Only show the tokens already expired or expiring within this number of days")
}
