package parsing

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

// This function parses a string that represents Modelica code and
// extracts all the libraries that it uses (along with their version
// numbers).  Sadly, one of Go's main weakness is its ability
// to built compiler compilers which makes creating a real Modelica
// parser problematic. :-(
func ParseUses(code string) (map[string]semver.Version, error) {
	// Empty result to return on error
	blank := map[string]semver.Version{}

	// Compile regexp to identify version strings
	ve, err := regexp.Compile(`version\s*=\s*"(.*)"`)
	if err != nil {
		return blank, fmt.Errorf("Unexpected error parsing 'version' regular expression: %v",
			err)
	}

	// Remove all tabs, spaces, line feeds and new lines
	compressed := strings.Replace(code, " ", "", -1)
	compressed = strings.Replace(compressed, "\n", "", -1)
	compressed = strings.Replace(compressed, "\r", "", -1)
	compressed = strings.Replace(compressed, "\t", "", -1)

	// Find first occurrence of uses(
	ustart := strings.Index(compressed, "uses(")
	if ustart == -1 {
		return blank, nil
	}

	// Take all text that appears after "uses("
	rem := compressed[ustart+5:]

	// Split what is left by ")"
	libs := strings.Split(rem, ")")

	// Initialize return value
	ret := map[string]semver.Version{}

	// Loop over all the split up chunks of text
	for _, lib := range libs {
		// If we find an empty chunk, it corresponds to a double ) which
		// marks the end of the uses annotation (so we are done)
		if lib == "" {
			break
		}

		// Split this chunk of text by "("s
		parts := strings.Split(lib, "(")

		// Trim any leading commas and this will be the library name
		libname := strings.TrimPrefix(parts[0], ",")

		// Look for version specifications contained in this chunk
		vers := ve.FindAllStringSubmatch(lib, 2)

		// If there are no versions, something is wrong with this annotation
		if len(vers) == 0 {
			return blank, fmt.Errorf("Unable to find version for library %s", libname)
		}
		// Similarly, if there are multiple versions, something is also wrong
		if len(vers) > 1 {
			return blank, fmt.Errorf("Multiple uses found for library %s???", libname)
		}

		// Record the (single) version we found
		nv, err := NormalizeVersion(vers[0][1])
		if err != nil {
			return blank, fmt.Errorf("Unable to normalize version for %s: %v", libname, err)
		}
		ret[libname] = nv
	}

	return ret, nil
}
