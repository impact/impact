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
	"encoding/json"

	"github.com/blang/semver"
)

type LocalLibrary struct {
	Name         string       `json:"name"`         // Name of library
	Path         string       `json:"path"`         // Path to library (relative to impact.json)
	IsFile       bool         `json:"isFile"`       // If the library is stored as a single file
	IssuesURL    string       `json:"issuesURL"`    // URL to issue tracker
	Dependencies []Dependency `json:"dependencies"` // Dependencies of this library
}

type Dependency struct {
	URI     string         `json:"uri"`     // Used to disambiguate libraries with the same name
	Name    string         `json:"name"`    // Name of library
	Version semver.Version `json:"version"` // Semantic version of library
}

type DirectoryInfo struct {
	Owner     string            `json:"owner"`     // Owner of this content
	Libraries []*LocalLibrary   `json:"libraries"` // Libraries defined here
	Alias     map[string]string `json:"alias"`     // key: library name, value: URI of library
}

func (di DirectoryInfo) JSON() string {
	bytes, err := json.MarshalIndent(di, "", "  ")
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ParseDirectoryInfo(s string) (DirectoryInfo, error) {
	ret := DirectoryInfo{}
	blank := DirectoryInfo{}

	err := json.Unmarshal([]byte(s), &ret)
	if err != nil {
		return blank, err
	}

	return ret, nil
}

func MakeDirectoryInfo() DirectoryInfo {
	return DirectoryInfo{
		Libraries: []*LocalLibrary{},
	}
}

func MakeLocalLibrary() LocalLibrary {
	return LocalLibrary{
		Dependencies: []Dependency{},
	}
}
