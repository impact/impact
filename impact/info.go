package main

import (
	"errors"
	"fmt"

	"github.com/xogeny/impact/index"

	"github.com/wsxiaoys/terminal/color"
)

type InfoCommand struct {
	Verbose    bool `short:"v" long:"verbose" description:"Turn on verbose output"`
	Positional struct {
		LibraryName string `description:"Library name"`
	} `positional-args:"true" required:"true"`
}

func (x *InfoCommand) Execute(args []string) error {
	if len(args) > 0 {
		fmt.Print("Ignoring (extra) unrecognized arguments: ")
		fmt.Println(args)
	}
	ind := index.DownloadIndex()

	libname := index.LibraryName(x.Positional.LibraryName)
	lib, ok := ind[libname]
	if !ok {
		return errors.New("Unable to locate library named '" + string(libname) + "'")
	}

	// TODO: Rewrite these as Printf's for God's sake!
	color.Println("@{!}Name: @{g}" + string(lib.Name))
	color.Println("  @{!}Description: @{c}" + string(lib.Description))
	color.Println("  @{!}Homepage: @{y}" + string(lib.Homepage))
	for vname, v := range lib.Versions {
		if x.Verbose {
			color.Println("    @{!}Version: " + string(vname))
			color.Println(fmt.Sprintf("      Major Version: %d", v.Version.Major))
			color.Println(fmt.Sprintf("      Minor Version: %d", v.Version.Minor))
			color.Println(fmt.Sprintf("      Patch Version: %d", v.Version.Patch))
			color.Println(fmt.Sprintf("      Pre-release tags: %d", v.Version.Pre))
			color.Println(fmt.Sprintf("      Build tags: %d", v.Version.Build))
			color.Println("      Canonical Version: " + v.Version.String())
		} else {
			color.Println("    @{!}Version: " + string(vname) + " (canonically " + v.Version.String() + ")")
		}
	}
	return nil
}
