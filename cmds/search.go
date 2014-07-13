package cmds

import "github.com/wsxiaoys/terminal/color"

/* Define a struct listing all command line options for 'search' */
type SearchCommand struct {
    URL bool `short:"u" long:"url" description:"Include homepage"`
}

var Search SearchCommand; // Instantiate option struct

/* This is the function called when the 'search' subcommand is executed */
func (x *SearchCommand) Execute(args []string) error {
	term := args[0];
	url := x.URL;

	index := buildIndex();
	for libname, lib := range(index) {
		if (lib.Matches(term)) {
			color.Println("@{g}"+string(libname)+":\n@{c}  - "+string(lib.Description))
			if (url) {
				color.Println("    @{y}"+lib.Homepage);
			}
		}
	}
	return nil
}
