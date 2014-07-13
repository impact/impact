package cmds

import "github.com/wsxiaoys/terminal/color"
import "errors"
import "strings"

/* Define a struct listing all command line options for 'install' */
type InstallCommand struct {
	Verbose bool `short:"v" login:"verbose" description:"Turn on verbose output"`
}

var Install InstallCommand; // Instantiate option struct

func (x *InstallCommand) Execute(args []string) error {
	if (len(args)==0) {
		return errors.New("No libraries requested for installation");
	}
	for _, arg := range(args) {
		var lib = arg
		var ver = "";
		parts := strings.Split(arg, "#");

		if (len(parts)==1) {
			// TODO: Determine latest version
		} else if (len(parts)==2) {
			lib = parts[0];
			ver = parts[1];
		} else if (len(parts)>2) {
			return errors.New("Invalid version specification: "+arg+
				" (must be libraryName#version)");
		}

		if (x.Verbose) {
			color.Println("@!Installing @{c}"+string(lib)+", version "+string(ver));
		}
	}
	return nil;
}
