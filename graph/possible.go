package graph

/*
 * This type represents remaining possible values for a given library
 */
type Possible map[LibraryName]*VersionList

/*
 * Create a clone of the map
 */
func (a Possible) Clone() Possible {
	clone := Possible{}
	for k, v := range a {
		clone[k] = v
	}
	return clone
}

/*
 * For each key in 'subset', only the values in that subset are still
 * possible.  So this method returns a new set of possible values which
 * represents either 1) the values provided in 'a' for keys of 'a' that
 * are not present in 'subset' or 2) the intersection of the values from
 * 'a' and 'subset' for keys they share.
 */
func (a Possible) Refine(subset Possible) Possible {
	ret := Possible{}

	for k, v := range a {
		v2, exists := subset[k]
		if !exists {
			ret[k] = v.Clone()
		} else {
			ret[k] = (*v).Intersection(*v2)
		}
		ret[k].ReverseSort()
	}
	return ret
}

/*
 * This method provides a list of libraries for which no more
 * possible versions exists.
 */
func (a Possible) Empty() []LibraryName {
	ret := []LibraryName{}
	for k, v := range a {
		if v.Len() == 0 {
			ret = append(ret, k)
		}
	}
	return ret
}
