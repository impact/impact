package graph

import (
	"fmt"

	"github.com/blang/semver"
)

/*
 * This struct is not exported.  It is used to represent a unique library
 * (i.e., name + version)
 */
type uniqueLibrary struct {
	name LibraryName
	ver  semver.Version
}

func (l1 uniqueLibrary) Equals(lib LibraryName, ver semver.Version) bool {
	return l1.name == lib && l1.ver.Compare(ver) == 0
}

/*
 * This is an edge in our (directed) dependency graph.  It indicates that `library`
 * depends on `dependsOn`.  Each is represented as a unique library (i.e., name + version)
 */
type dependency struct {
	library   uniqueLibrary
	dependsOn uniqueLibrary
}

func (d dependency) String() string {
	return fmt.Sprintf("%s:%s -> %s:%s", d.library.name,
		d.library.ver.String(), d.dependsOn.name,
		d.dependsOn.ver.String())
}

func (d1 dependency) Equals(d2 dependency) bool {
	return d1.library.name == d2.library.name &&
		d1.library.ver.Compare(d2.library.ver) == 0 &&
		d1.dependsOn.name == d2.dependsOn.name &&
		d1.dependsOn.ver.Compare(d2.dependsOn.ver) == 0
}
