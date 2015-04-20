package main

import (
	"errors"
	"fmt"

	"github.com/wsxiaoys/terminal/color"

	"github.com/xogeny/impact/graph"
	"github.com/xogeny/impact/index"
	"github.com/xogeny/impact/install"
)

/* Define a struct listing all command line options for 'install' */
type InstallCommand struct {
	Verbose bool `short:"v" long:"verbose" description:"Turn on verbose output"`
	DryRun  bool `short:"d" long:"dryrun" description:"Resolve dependencies but don't install"`
}

func (x InstallCommand) Execute(args []string) error {
	//Check to make sure we have something to install
	if len(args) == 0 {
		return errors.New("No libraries requested for installation")
	}

	// Load index
	ind, err := index.LoadIndex()
	if err != nil {
		return fmt.Errorf("Error loading indices: %v", err)
	}

	// Build dependency graph from index
	resolver, err := ind.BuildGraph(x.Verbose)
	if err != nil {
		return fmt.Errorf("Error building dependency graph: %v", err)
	}

	// State root dependencies
	libnames := []graph.LibraryName{}
	for _, n := range args {
		libnames = append(libnames, graph.LibraryName(n))
	}

	// Resolve dependencies
	solution, err := resolver.Resolve(libnames...)
	if err != nil {
		return fmt.Errorf("Error resolving dependencies for %v: %v", libnames, err)
	}

	// Install dependencies
	color.Printf("@{y}Installing...\n")
	for name, version := range solution {
		color.Printf("  @{g}Library: @{!g}%s\n", name)
		color.Printf("    @{g}Required version: @{!g}%v\n", version)
		lv, err := ind.Find(string(name), version)
		if err != nil {
			return fmt.Errorf("Couldn't find version %v of library %s (this should not happen)",
				version, name)
		}
		if !x.DryRun {
			install.Install(string(name), lv, ind, ".", x.Verbose)
		}
	}

	return nil
}
