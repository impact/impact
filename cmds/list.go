package cmds

import "fmt"

import "github.com/xogeny/impact/utils"
import "github.com/wsxiaoys/terminal/color"

/* Define a struct listing all command line options for 'search' */
type ListCommand struct {
	URL bool `short:"u" long:"url" description:"Include homepage"`
}

/* This is the function called when the 'search' subcommand is executed */
func (x *ListCommand) Execute(args []string) error {
	if len(args) > 0 {
		fmt.Print("Ignoring (extra) unrecognized arguments: ")
		fmt.Println(args)
	}
	url := x.URL

	index := utils.DownloadIndex()
	for libname, lib := range index {
		color.Println("@{g}" + string(libname) + ":\n@{c}  - " + string(lib.Description))
		if url {
			color.Println("    @{y}" + lib.Homepage)
		}
	}
	return nil
}
