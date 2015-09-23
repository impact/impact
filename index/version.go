package index

import (
	"github.com/blang/semver"

	"github.com/impact/impact/recorder"
)

type VersionDetails struct {
	Version semver.Version `json:"version"`
	Tarball string         `json:"tarball_url"`
	Zipball string         `json:"zipball_url"`

	// This indicates where (within an archive) the library can be found:
	Path string `json:"path"`
	// This indicates whether the specified path is to a file or directory:
	IsFile bool `json:"isfile"`

	Dependencies []Dependency `json:"dependencies"`
	Sha          string       `json:"sha"`
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

func (v *VersionDetails) SetPath(path string, file bool) {
	v.Path = path
	v.IsFile = file
}

func (v *VersionDetails) AddDependency(library string, version semver.Version) {
	v.Dependencies = append(v.Dependencies, Dependency{
		Name:    library,
		Version: version.String(),
	})
}

var _ recorder.VersionRecorder = (*VersionDetails)(nil)
