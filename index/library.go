package index

import (
	"strings"

	"github.com/blang/semver"

	"github.com/impact/impact/recorder"
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
	Repository string `json:"repository_uri"`
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

func (lib Library) Matches(term string) bool {
	match := false

	t := strings.ToLower(term)
	lname := strings.ToLower(lib.Name)
	ldesc := strings.ToLower(lib.Description)

	if strings.Contains(lname, t) {
		match = true
	}
	if strings.Contains(ldesc, t) {
		match = true
	}
	return match
}

var _ recorder.LibraryRecorder = (*Library)(nil)
