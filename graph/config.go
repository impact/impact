package graph

import (
	"github.com/blang/semver"
)

/*
 * This type represents a specific configuration of libraries.  This is used to represent
 * the resolution of dependencies.
 */
type Configuration map[LibraryName]semver.Version

func (conf Configuration) Clone() Configuration {
	clone := Configuration{}
	for k, v := range conf {
		clone[k] = v
	}
	return clone
}
