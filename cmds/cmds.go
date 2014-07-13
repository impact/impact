package cmds

import "fmt"
import "xogeny/gimpact/utils"
import "github.com/wsxiaoys/terminal/color"
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
			color.Println("@{g}"+string(libname)+":\n@{c}  - "+string(lib.Description))
		}
	}
}
