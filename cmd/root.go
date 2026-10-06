package cmd

import (
	"embed"
	"opsi/config"
	"opsi/helpers"
	"opsi/helpers/ui"
	git "opsi/scopes/gitlab"
	host "opsi/scopes/hosts"
	op "opsi/scopes/onepassword"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainConfig config.Config

var gitlab git.Gitlab
var onepassword op.OnePassword
var hosts host.Host

var ConfigTemplate embed.FS

// Version of the app provided
// in build phase
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "opsi",
	Version: Version,
	Short:   "All-in-one CLI for Beliven Ops daily usage",
	Long: `All-in-one CLI for Beliven Ops daily usage.

Commands are grouped by scope (gitlab, 1password, hosts) and then by the
thing they act on: opsi <scope> <entity> <verb>, for example
opsi gitlab project create.

The configuration lives in ~/.config/opsi/config.yml and is created on the
first run. Add --help to any command to see its flags and examples.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Cobra would print the error and the whole usage on its own: the error
	// goes through ui, followed by a pointer to the help of the command
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	command, err := rootCmd.ExecuteC()
	if err != nil {
		ui.Error("%s", err.Error())
		ui.Muted("Run '%s --help' for usage", command.CommandPath())
		os.Exit(1)
	}
}

// showHelp is the run of the commands that only group other commands
func showHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func initConfig() {
	// Checked here rather than in init so that importing the package, and every
	// command that does not shell out to op, does not depend on it
	if helpers.Which("op") == "" {
		ui.Error("Missing op executable")
		ui.Muted("Please follow the instructions for install the binaries here:")
		ui.Muted("https://developer.1password.com/docs/cli/get-started")
		os.Exit(1)
	}

	// Don't forget to read config either from cfgFile or from home directory!
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		ui.Fatal(err)
	}

	var configFolder = "/.config/opsi/"
	var configName = "config"

	helpers.ConfigInit(ConfigTemplate, "/.config/opsi/config.yml")

	// Check if the user just use exist
	if _, err := os.Stat(home + "/.opsi.yml"); err == nil {
		configFolder = ""
		configName = ".opsi"
	}

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home + configFolder)
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")

	// Read config
	err = viper.ReadInConfig()
	if err != nil {
		ui.Error("Config file error: %s", err.Error())
		os.Exit(0)
	}

	if err := viper.Unmarshal(&mainConfig); err != nil {
		ui.Fatal(err)
	}

	gitlab = git.NewGitlab(
		mainConfig.Gitlab.ApiURL,
		mainConfig.Gitlab.Token,
		mainConfig.Gitlab.Mirror,
		mainConfig.Gitlab.Exclusions,
	)

	onepassword = op.NewOnePassword(mainConfig.OnePassword.Address)

	hosts = host.NewHosts()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Headings of the help in bold, only on a terminal
	cobra.AddTemplateFunc("bold", ui.Bold)
	cobra.AddTemplateFunc("dim", ui.Dim)

	template := rootCmd.UsageTemplate()
	for _, heading := range []string{"Usage:", "Aliases:", "Examples:", "Available Commands:", "Additional Commands:", "Flags:", "Global Flags:", "Additional help topics:"} {
		template = strings.Replace(template, heading, `{{bold "`+heading+`"}}`, 1)
	}
	template = strings.Replace(template, "{{.Title}}", "{{bold .Title}}", 1)
	template = strings.Replace(template, `Use "{{.CommandPath}} [command] --help" for more information about a command.`,
		`{{dim (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}`, 1)
	template = strings.Replace(template, "{{if .Runnable}}", "{{if and .Runnable (not .HasAvailableSubCommands)}}", 1)
	rootCmd.SetUsageTemplate(template)
}
