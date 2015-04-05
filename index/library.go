package index

import (
	"strings"

	"github.com/blang/semver"

	"github.com/xogeny/impact/recorder"
)

type Library struct {
	Owner       string                     `json:"owner"`
	Email       string                     `json:"email"`
	Name        string                     `json:"name"`
	Homepage    string                     `json:"homepage"`
	Description string                     `json:"description"`
	Stars       int                        `json:"stars"`
	Versions    map[string]*VersionDetails `json:"versions"`
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

func (lib *Library) AddVersion(v semver.Version) recorder.VersionRecorder {
	details := NewVersionDetails(v)
	lib.Versions[v.String()] = details
	return details
}

func NewLibrary(owner string, name string) *Library {
	return &Library{
		Owner:    owner,
		Name:     name,
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
