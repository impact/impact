package graph

import "github.com/blang/semver"

type Resolver interface {
	Contains(LibraryName, semver.Version) bool
	AddLibrary(LibraryName, semver.Version)
	AddDependency(LibraryName, semver.Version, LibraryName, semver.Version) error
	Resolve(...LibraryName) (Configuration, error)
}
