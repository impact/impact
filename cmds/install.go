package cmds

import "github.com/wsxiaoys/terminal/color"
import "errors"

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
		if (x.Verbose) {
			color.Println("@!Installing @{c}"+string(arg));
		}
	}
	return nil;
}
