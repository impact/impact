package index

import (
	"sort"
)

type LibList []*Library

func (l LibList) Len() int {
	return len(l)
}

func (l LibList) Swap(i int, j int) {
	x := l[i]
	l[i] = l[j]
	l[j] = x
}

func (l LibList) Less(i int, j int) bool {
	if i == 0 {
		return true
	}
	if j == 0 {
		return false
	}
	return l[i].Stars > l[j].Stars
}

type GroupedIndex struct {
	Libraries map[string]LibList
}

func makeGroupedIndex() GroupedIndex {
	return GroupedIndex{
		Libraries: map[string]LibList{},
	}
}

func (i Index) Group(disamb map[string]string) GroupedIndex {
	ret := makeGroupedIndex()
	for _, lib := range i.Libraries {
		list, exists := ret.Libraries[lib.Name]
		if exists {
			// Check to see if this library is disambiguated
			curi, selection := disamb[lib.Name]
			if selection && curi == lib.URI {
				// If this is the selected version, put it at the
				// front of the list
				ret.Libraries[lib.Name] = append(LibList{lib}, list...)
			} else {
				// Otherwise, just add it to the list
				ret.Libraries[lib.Name] = append(list, lib)
			}
		} else {
			ret.Libraries[lib.Name] = LibList{lib}
		}
	}
	return ret
}

// This sorts **all but the first** match by rating.  This is because the
// first match is the one that is first in the install path so it should
// always be listed first.  Among those not in the search path, this
// rearranges them by rating.
func (g GroupedIndex) SortByRating() GroupedIndex {
	ret := makeGroupedIndex()

	for libname, libs := range g.Libraries {
		sort.Sort(libs)
		ret.Libraries[libname] = libs
	}
	return ret
}

// Return and index that has no "extra" libraries in it, only the ones that
// can possibly be chosen.
func (g GroupedIndex) Selected() *Index {
	ret := NewIndex()

	for _, libs := range g.Libraries {
		ret.Libraries = append(ret.Libraries, libs[0])
	}

	return ret
}
