package cmds

import "fmt"
import "errors"

import "github.com/wsxiaoys/terminal/color"
import "github.com/xogeny/gimpact/utils"

type InfoCommand struct {
	Verbose    bool `short:"v" login:"verbose" description:"Turn on verbose output"`
	Positional struct {
		LibraryName string `description:"Library name"`
	} `positional-args:"true" required:"true"`
}

func (x *InfoCommand) Execute(args []string) error {
	if len(args) > 0 {
		fmt.Print("Ignoring (extra) unrecognized arguments: ")
		fmt.Println(args)
	}
	index := utils.DownloadIndex()

	libname := utils.LibraryName(x.Positional.LibraryName)
	lib, ok := index[libname]
	if !ok {
		return errors.New("Unable to locate library named '" + string(libname) + "'")
	}
	color.Println("@{!}Name: @{g}" + string(lib.Name))
	color.Println("  @{!}Description: @{c}" + string(lib.Description))
	color.Println("  @{!}Homepage: @{y}" + string(lib.Homepage))
	for vname, v := range lib.Versions {
		if x.Verbose {
			color.Println("    @{!}Version: " + string(vname))
			color.Println(fmt.Sprintf("      Major Version: %d", v.Major))
			color.Println(fmt.Sprintf("      Minor Version: %d", v.Minor))
			color.Println(fmt.Sprintf("      Patch Version: %d", v.Patch))
			color.Println("      Canonical Version: " + string(v.Version))
		} else {
			color.Println("    @{!}Version: " + string(vname) + " (canonically " + string(v.Version) + ")")
		}
	}
	return nil
}
