package cmd

import (
	"embed"
	"fmt"
	"opsi/config"
	"opsi/helpers"
	git "opsi/scopes/gitlab"
	host "opsi/scopes/hosts"
	op "opsi/scopes/onepassword"
	"os"

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
	Short:   "The root command",
	Long:    `The root command`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Println("Config file error", err.Error())
		os.Exit(0)
	}

	if err := viper.Unmarshal(&mainConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gitlab = git.NewGitlab(mainConfig.Gitlab.ApiURL, mainConfig.Gitlab.Token, mainConfig.Gitlab.GroupID)
	onepassword = op.NewOnePassword(mainConfig.OnePassword.Address)
	hosts = host.NewHosts()
}

func init() {
	if helpers.Which("op") == "" {
		fmt.Println("Missing op executable")
		fmt.Println("Please follow the instructions for install the binaries here:")
		fmt.Print("\nhttps://developer.1password.com/docs/cli/get-started\n\n")
		os.Exit(1)
	}

	cobra.OnInitialize(initConfig)
}
