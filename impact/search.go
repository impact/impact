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
	Ratings bool `short:"r" long:"ratings" description:"Show ratings"`
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
	ratings := x.Ratings

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

	g := ind.Group(settings.Choices).SortByRating()

	for _, libs := range g.Libraries {
		multiple := len(libs) > 0 && x.All
		if multiple {
			url = true     // Always show URLs with multiple matches
			ratings = true // Always show ratings with multiple matches
		}
		for index, lib := range libs {
			prefix := ""
			if index > 0 && !x.All {
				continue
			}
			if index > 0 {
				prefix = "  "
			}
			// TODO: Should have some way of indicating the "chosen" one (i.e., what
			// would be installed).
			if lib.Matches(term) {
				rstr := ""
				if ratings {
					if lib.Stars > 0 {
						rstr = fmt.Sprintf(" [%d stars]", lib.Stars)
					} else {
						rstr = fmt.Sprintf(" [no rating]", lib.Stars)
					}
				}
				if !x.All && len(libs) > 0 {
					rstr = fmt.Sprintf("%s (use -a flag to see %d other libraries named %s)",
						rstr, len(libs)-1, lib.Name)
				}

				if multiple && index == 0 {
					color.Printf("@{y}First match in your search path:\n")
				}
				if multiple && index == 1 {
					color.Printf("\n@{y}%sAlternatives:\n", prefix)
				}
				if url {
					color.Printf("%s@{!}%s%s @{y}(%s):\n%s@{c} - %s\n", prefix, lib.Name,
						rstr, lib.URI, prefix, lib.Description)
				} else {
					color.Printf("%s@{!}%s%s:\n%s@{c} - %s\n", prefix, lib.Name, rstr,
						prefix, lib.Description)
				}
			}
		}
	}
	return nil
}
