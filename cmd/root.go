/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"embed"
	"fmt"
	"opsi/config"
	"opsi/helpers"
	"opsi/scopes"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainConfig config.Config

var gitlab scopes.Gitlab
var postmark scopes.Postmark
var onepassword scopes.OnePassword
var hosts scopes.Hosts

var Scripts embed.FS

var ConfigTemplate embed.FS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opsi",
	Short: "The root command",
	Long:  `The root command`,
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

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".opsi")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		// Initialize and exit
		err = helpers.ConfigInit(ConfigTemplate)
		if err != nil {
			fmt.Println("Cannot Initialize config")
			os.Exit(1)
		}

		fmt.Println("Config create successfully")

		os.Exit(0)
	}

	if err := viper.Unmarshal(&mainConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gitlab = scopes.NewGitlab(mainConfig.Gitlab.ApiURL, mainConfig.Gitlab.Token, mainConfig.Gitlab.GroupID)
	postmark = scopes.NewPostmark(mainConfig.Postmark.ApiURL, mainConfig.Postmark.Token, mainConfig.Postmark.SlackWebhook)
	onepassword = scopes.NewOnePassword(mainConfig.OnePassword.Address)
	hosts = scopes.NewHosts()
}

func init() {
	cobra.OnInitialize(initConfig)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.opsi.yaml)")
}
