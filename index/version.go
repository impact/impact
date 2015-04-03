package index

import (
	"github.com/blang/semver"
)

type VersionDetails struct {
	Version      semver.Version `json:"version"`
	Tarball      string         `json:"tarball_url"`
	Zipball      string         `json:"zipball_url"`
	Path         string         `json:"path"`
	Dependencies []Dependency   `json:"dependencies"`
	Sha          string         `json:"sha"`
}
