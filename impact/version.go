package main

import (
	"github.com/wsxiaoys/terminal/color"
)

var version = "0.7.0"

type VersionCommand struct {
	Verbose bool `short:"v" long:"verbose" description:"Turn on verbose output"`
}

func (x VersionCommand) Execute(args []string) error {
	color.Printf("@{g}Impact version @{!g}%s\n", version)
	return nil
}
