package crawl

import (
	"github.com/blang/semver"
)

type Recorder interface {
	AddLibrary(name string) LibraryRecorder
}

type LibraryRecorder interface {
	SetStars(int)
	AddVersion(v semver.Version) VersionRecorder
}

type VersionRecorder interface {
	SetHash(hash string)
	SetTarballURL(url string)
	SetZipballURL(url string)
	AddDependency(library string, version semver.Version)
}
