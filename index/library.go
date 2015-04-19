package index

import (
	"strings"

	"github.com/blang/semver"

	"github.com/xogeny/impact/recorder"
)

type Library struct {
	// Library name
	Name string `json:"name"`
	// Canonical URI for this library (used for disambiguation)
	URI string `json:"uri"`

	// Versions of this library that are available
	Versions map[string]*VersionDetails `json:"versions"`

	// Owner of library
	OwnerURI string `json:"owner_uri"`
	// Email for library contact
	Email string `json:"email"`
	// Web site
	Homepage string `json:"homepage"`
	// Repository
	Repository string `json:"repository"`
	// Repository format
	Format string `json:"repository_format"`
	// Textual description
	Description string `json:"description"`
	// Stars (if applicable, otherwise -1)
	Stars int `json:"stars"`
}

func (lib *Library) SetEmail(email string) {
	lib.Email = email
}

func (lib *Library) SetStars(stars int) {
	lib.Stars = stars
}

func (lib *Library) SetDescription(desc string) {
	lib.Description = desc
}

func (lib *Library) SetHomepage(url string) {
	lib.Homepage = url
}

func (lib *Library) SetRepository(url string, format string) {
	lib.Repository = url
	lib.Format = format
}

func (lib *Library) AddVersion(v semver.Version) recorder.VersionRecorder {
	details := NewVersionDetails(v)
	lib.Versions[v.String()] = details
	return details
}

func NewLibrary(name string, uri string, owner_uri string) *Library {
	return &Library{
		Name:     name,
		URI:      uri,
		OwnerURI: owner_uri,
		Stars:    -1,
		Versions: map[string]*VersionDetails{},
	}
}

var _ recorder.LibraryRecorder = (*Library)(nil)

/*
func (lib Library) Latest() (ver VersionDetails, err error) {
	var first = true
	if len(lib.Versions) == 0 {
		err = EmptyLibraryError{Name: lib.Name}
		return
	}
	for _, v := range lib.Versions {
		if first {
			ver = v
			first = false
		}
		if !first && v.Version.GT(ver.Version) {
			ver = v
		}
	}
	return
}
*/

func (lib Library) Matches(term string) bool {
	var match = false
	if strings.Contains(string(lib.Name), term) {
		match = true
	}
	if strings.Contains(string(lib.Description), term) {
		match = true
	}
	return match
}

type Libraries map[string]VersionDetails

/*
func (libs *Libraries) Merge(olibs Libraries) error {
	for ln, lv := range olibs {
		ev, exists := (*libs)[ln]
		if exists {
			if lv.Version.Equals(ev.Version) {
				return VersionConflictError{Name: ln, Existing: ev, Additional: lv}
			}
		} else {
			(*libs)[ln] = lv
		}
	}
	return nil
}
*/
