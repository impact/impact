package main

import "os"

import "github.com/xogeny/impact/cmds"

import "github.com/jessevdk/go-flags"

func main() {
	var options struct{} // No common flags

	parser := flags.NewParser(&options, flags.Default)

	parser.AddCommand("search",
		"Search for libraries associated with specific terms",
		"Search for libraries associated with specific terms",
		&cmds.SearchCommand{})

	parser.AddCommand("install",
		"Install named libraries",
		"Install named libraries",
		&cmds.InstallCommand{})

	parser.AddCommand("info",
		"Print information about a specific library",
		"Print information about a specific library",
		&cmds.InfoCommand{})

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
