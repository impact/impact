package index

import (
	"github.com/blang/semver"
	"github.com/xogeny/impact/recorder"
)

type versionRecorder struct {
}

func (vr versionRecorder) SetHash(hash string)                                  {}
func (vr versionRecorder) SetTarballURL(url string)                             {}
func (vr versionRecorder) SetZipballURL(url string)                             {}
func (vr versionRecorder) AddDependency(library string, version semver.Version) {}

func makeVersionRecorder() versionRecorder {
	return versionRecorder{}
}

var _ recorder.VersionRecorder = (*versionRecorder)(nil)
