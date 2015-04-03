package index

import (
	"strings"
)

type Library struct {
	Name        LibraryName
	Homepage    string                           `json:"homepage"`
	Description string                           `json:"description"`
	Versions    map[VersionString]VersionDetails `json:"versions"`
}

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

type Libraries map[LibraryName]VersionDetails

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
