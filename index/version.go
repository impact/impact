package index

import (
	"github.com/blang/semver"

	"github.com/xogeny/impact/recorder"
)

type VersionDetails struct {
	Version      semver.Version `json:"version"`
	Tarball      string         `json:"tarball_url"`
	Zipball      string         `json:"zipball_url"`
	Path         string         `json:"path"`
	Dependencies []Dependency   `json:"dependencies"`
	Sha          string         `json:"sha"`
}

func NewVersionDetails(v semver.Version) *VersionDetails {
	return &VersionDetails{
		Version:      v,
		Dependencies: []Dependency{},
	}
}

func (v *VersionDetails) SetHash(hash string) {
	v.Sha = hash
}

func (v *VersionDetails) SetTarballURL(url string) {
	v.Tarball = url
}

func (v *VersionDetails) SetZipballURL(url string) {
	v.Zipball = url
}

func (v *VersionDetails) AddDependency(library string, version semver.Version) {
	v.Dependencies = append(v.Dependencies, Dependency{
		Name:    library,
		Version: version.String(),
	})
}

var _ recorder.VersionRecorder = (*VersionDetails)(nil)
