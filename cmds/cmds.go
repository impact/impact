package cmds

import "fmt"
import "xogeny/gimpact/utils"
import "strings"

func Search(index utils.Index, term string, verbose bool) {
	if (verbose) {
		fmt.Println("Searching for "+term);
	}

	for libname, lib := range(index) {
		var match = false;
		if (strings.Contains(string(libname), term)) {
			match = true;
		}
		if (strings.Contains(string(lib.Description), term)) {
			match = true;
		}
		if (match) {
			fmt.Println(string(libname)+": "+string(lib.Description));
		}
	}
}
