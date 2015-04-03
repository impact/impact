package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct{} // No common flags

	parser := flags.NewParser(&options, flags.Default)

	parser.AddCommand("search",
		"Search for libraries associated with specific terms",
		"Search for libraries associated with specific terms",
		&SearchCommand{})

	parser.AddCommand("list",
		"List all available libraries",
		"List all available libraries",
		&cmds.ListCommand{})

	parser.AddCommand("install",
		"Install named libraries",
		"Install named libraries",
		&InstallCommand{})

	parser.AddCommand("info",
		"Print information about a specific library",
		"Print information about a specific library",
		&InfoCommand{})

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
