package index

type GroupedIndex struct {
	Libraries map[string][]*Library
}

func makeGroupedIndex() GroupedIndex {
	return GroupedIndex{
		Libraries: map[string][]*Library{},
	}
}

func (i Index) GroupByOrder(disamb map[string]string) GroupedIndex {
	ret := makeGroupedIndex()
	for _, lib := range i.Libraries {
		list, exists := ret.Libraries[lib.Name]
		if exists {
			// Check to see if this library is disambiguated
			curi, selection := disamb[lib.Name]
			if selection && curi == lib.URI {
				// If this is the selected version, put it at the
				// front of the list
				ret.Libraries[lib.Name] = append([]*Library{lib}, list...)
			} else {
				// Otherwise, just add it to the list
				ret.Libraries[lib.Name] = append(list, lib)
			}
		} else {
			ret.Libraries[lib.Name] = []*Library{lib}
		}
	}
	return ret
}

func (i Index) GroupByRating() GroupedIndex {
	return makeGroupedIndex()
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
