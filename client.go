package main

import "os"

import "xogeny/gimpact/cmds"

import "github.com/jessevdk/go-flags"

func main() {
	var options struct {}; // No common flags
	
	parser := flags.NewParser(&options, flags.Default)

	parser.AddCommand("search",
		"Search for libraries associated with specific terms",
		"Search for libraries associated with specific terms",
		&cmds.Search)

	parser.AddCommand("install",
		"Install named libraries",
		"Install named libraries",
		&cmds.Install)

	parser.AddCommand("info",
		"Print information about a specific library",
		"Print information about a specific library",
		&cmds.Info)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
