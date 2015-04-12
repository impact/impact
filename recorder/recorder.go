package recorder

import (
	"github.com/blang/semver"
)

type Recorder interface {
	// Create library if it doesn't already exist.  Otherwise, return
	// recorder for existing library
	// The first argument, 'name', is the name of the library (in Modelica terms).
	// The second argument, 'uri', is a URI indicating precisely what library this
	//   is.  This helps address things like forks, etc.
	// The third argument is the owner (as a URI)
	GetLibrary(name string, uri string, owner_uri string) LibraryRecorder
}

type LibraryRecorder interface {
	SetDescription(desc string)
	SetHomepage(url string)
	SetStars(int)
	SetEmail(string)
	AddVersion(v semver.Version) VersionRecorder
}

type VersionRecorder interface {
	SetHash(hash string)
	SetTarballURL(url string)
	SetZipballURL(url string)
	SetPath(path string, file bool)
	AddDependency(library string, version semver.Version)
}
