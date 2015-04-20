package index

type GroupedIndex struct {
	Libraries map[string][]*Library
}

func makeGroupedIndex() GroupedIndex {
	return GroupedIndex{
		Libraries: map[string][]*Library{},
	}
}

func (i Index) GroupByOrder() GroupedIndex {
	ret := makeGroupedIndex()
	for _, lib := range i.Libraries {
		list, exists := ret.Libraries[lib.Name]
		if exists {
			ret.Libraries[lib.Name] = append(list, lib)
		} else {
			ret.Libraries[lib.Name] = []*Library{lib}
		}
	}
	return ret
}

func (i Index) GroupByRating() GroupedIndex {
	return makeGroupedIndex()
}
