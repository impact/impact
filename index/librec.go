package index

import (
	"github.com/blang/semver"
	"github.com/xogeny/impact/crawl"
)

type libraryRecorder struct {
	lib *Library
}

func (lr libraryRecorder) SetStars(int) {}

func (lr libraryRecorder) AddVersion(v semver.Version) crawl.VersionRecorder {
	return makeVersionRecorder()
}

func makeLibraryRecorder(lib *Library) libraryRecorder {
	return libraryRecorder{
		lib: lib,
	}
}

var _ crawl.LibraryRecorder = (*libraryRecorder)(nil)
