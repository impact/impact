package main

import (
	"fmt"

	"github.com/wsxiaoys/terminal/color"

	"github.com/xogeny/impact/config"
	"github.com/xogeny/impact/index"
)

/* Define a struct listing all command line options for 'search' */
type SearchCommand struct {
	Positional struct {
		Term string `description:"Search term"`
	} `positional-args:"true" required:"true"`
	All     bool `short:"a" long:"all" description:"Include all libraries in index"`
	URL     bool `short:"u" long:"url" description:"Include homepage"`
	Verbose bool `short:"v" long:"verbose" description:"Turn on verbose output"`
}

/* This is the function called when the 'search' subcommand is executed */
func (x *SearchCommand) Execute(args []string) error {
	if len(args) > 0 {
		fmt.Print("Ignoring (extra) unrecognized arguments: ")
		fmt.Println(args)
	}
	term := x.Positional.Term
	url := x.URL

	// Load user settings
	settings, err := config.ReadSettings()
	if err != nil {
		return fmt.Errorf("Error reading settings: %v", err)
	}

	// Load index
	ind, err := index.LoadIndex(x.Verbose)
	if err != nil {
		return fmt.Errorf("Error loading indices: %v", err)
	}

	if x.All {
		return fmt.Errorf("The --all flag in search is currently not implemented")
	} else {
		ind = ind.Reduce(settings.Choices)
	}

	for _, lib := range ind.Libraries {
		if lib.Matches(term) {
			if url {
				color.Printf("@{!}%s @{y}(%s):\n@{c} - %s\n", lib.Name, lib.Homepage,
					lib.Description)
			} else {
				color.Printf("@{!}%s:\n@{c} - %s\n", lib.Name, lib.Description)
			}
		}
	}
	return nil
}
