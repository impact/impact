// This package represents information we need to know about libraries
// stored somewhere (e.g., GitHub).  This is information we need to
// build the index.
//
// HOWEVER, this information is not required to be provided.  We will use
// a collection of heuristics to construct this information for any given
// directory.  BUT, in complex cases where heuristics are not possible
// or library developers wish to follow a different set of conventions,
// they have the option of explicitly creating an 'impact.json' file
// to provide this information.
//
// Note: Where information is missing (empty strings, empty arrays), an
// attempt should be made to infer the information.  Only when it is
// explicitly provided should we avoid the use of heuristics to fill in
// the information
package dirinfo

import (
	"github.com/blang/semver"
)

type LocalLibrary struct {
	Owner        string       // Owner of library (to distinguish libs with same name)
	Name         string       // Name of library
	Path         string       // Path to library (relative to directory/impact.json file)
	IsFile       bool         // If the library is stored as a single file
	IssuesURL    string       // URL to issue tracker
	Dependencies []Dependency // Dependencies of this library
}

type Dependency struct {
	Owner         string         // Owner of library
	Name          string         // Name of library
	Version       semver.Version // Semantic version of library
	VersionString string         // Version string (in case not semver)
}

type DirectoryInfo struct {
	Owner     string         // Owner of this content
	Libraries []LocalLibrary // Libraries defined here
}
