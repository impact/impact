package cmds

//import "github.com/wsxiaoys/terminal/color"
import "fmt"

/* Define a struct listing all command line options for 'install' */
type InstallCommand struct {
	Verbose bool `short:"v" login:"verbose" description:"Turn on verbose output"`
}

var Install InstallCommand; // Instantiate option struct

func (x *InstallCommand) Execute(args []string) error {
	fmt.Println("verbose:");
	fmt.Println(x.Verbose);
	return nil;
}
