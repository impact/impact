package index

import (
	"strings"
)

type Library struct {
	Name        LibraryName
	Homepage    string                    `json:"homepage"`
	Description string                    `json:"description"`
	Versions    map[VersionString]Version `json:"versions"`
}

func (lib Library) Latest() (ver Version, err error) {
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
		if !first && v.GreaterThan(ver) {
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

type Libraries map[LibraryName]Version

func (libs *Libraries) Merge(olibs Libraries) error {
	for ln, lv := range olibs {
		ev, exists := (*libs)[ln]
		if exists {
			if lv.Equals(ev) {
				return VersionConflictError{Name: ln, Existing: ev, Additional: lv}
			}
		} else {
			(*libs)[ln] = lv
		}
	}
	return nil
}
