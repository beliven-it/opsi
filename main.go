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

func main() {
	cmd.Scripts = scripts

	cmd.Execute()
}
