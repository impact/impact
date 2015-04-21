package main

import (
	"github.com/wsxiaoys/terminal/color"

	"github.com/xogeny/impact/config"
)

var version = "0.9.0-dev"

type VersionCommand struct {
	Verbose bool `short:"v" long:"verbose" description:"Turn on verbose output"`
}

func (x VersionCommand) Execute(args []string) error {
	color.Printf("@{g}Impact version @{!g}%s\n", version)
	color.Printf("  Settings file: %s\n", config.SettingsFile())
	settings, err := config.ReadSettings()
	if err != nil {
		color.Printf("    @{!r}%s\n", err.Error())
	} else {
		color.Printf("@{g}%s\n", settings.List(""))
	}
	return nil
}
