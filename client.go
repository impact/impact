package main

import "fmt"
import "xogeny/gimpact/utils"
import "net/http"
import "flag"
import "xogeny/gimpact/cmds"

func main() {
	var master = "http://impact.modelica.org/impact_data.json";

	// Option variables
	var verbose bool;

	// Option definitions
	flag.BoolVar(&verbose, "verbose", false, "Turn on verbose output")

	// Parse options
	flag.Parse();

	args := flag.Args();
	
	if (len(args)<1) {
		fmt.Println("No command specified, exiting.");
		return;
	}

	resp, err := http.Get(master)
	if err != nil {
		fmt.Println("Unable to locate index file at "+master);
		return;
	}
	defer resp.Body.Close()

	index := utils.Index{};

	err = index.BuildIndex(resp.Body);
	if (err!=nil) {
		fmt.Println("Error reading index: "+err.Error());
	}

	cmd := args[0];

	switch cmd {
	case "search":
		if (len(args)<2) {
			fmt.Println("No search term given");
			return;
		} else if (len(args)>2) {
			fmt.Println("Extra data after search term");
			return;
		}
		cmds.Search(index, args[1], verbose);
	default:
		fmt.Println("Unrecognized command "+cmd);
		return;
	}

	if (verbose) {
		fmt.Println("Parsing index at "+master);
	}

	// Add command parsing here
}
