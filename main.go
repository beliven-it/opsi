/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"
	"opsi/cmd"
)

//go:embed scripts/*
var scripts embed.FS

//go:embed config/template.yml
var configTemplate embed.FS

func main() {
	cmd.Scripts = scripts
	cmd.ConfigTemplate = configTemplate

	cmd.Execute()
}
