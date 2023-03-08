/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"
	"opsi/cmd"
)

//go:embed config/template.yml
var configTemplate embed.FS

func main() {
	cmd.ConfigTemplate = configTemplate

	cmd.Execute()
}
