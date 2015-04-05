package recorder

import (
	"github.com/blang/semver"
)

type Recorder interface {
	AddLibrary(owner string, name string) LibraryRecorder
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
	AddDependency(library string, version semver.Version)
}
