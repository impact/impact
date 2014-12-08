package deps

type LibraryName string

type uniqueLibrary struct {
	name LibraryName
	ver  SemanticVersion
}

type dependency struct {
	library   uniqueLibrary
	dependsOn uniqueLibrary
}

type LibraryIndex struct {
	libraries []dependency
}

func MakeLibraryIndex() LibraryIndex {
	return LibraryIndex{
		libraries: []dependency{},
	}
}

func (index *LibraryIndex) AddDependency(lib LibraryName, libver SemanticVersion,
	dep LibraryName, depver SemanticVersion) {

	library := uniqueLibrary{name: lib, ver: libver}
	dependsOn := uniqueLibrary{name: dep, ver: depver}
	index.libraries = append(index.libraries, dependency{library: library, dependsOn: dependsOn})
}

func addVersion(vers []SemanticVersion, v SemanticVersion) []SemanticVersion {
	l := len(vers)
	if l == 0 {
		return []SemanticVersion{v}
	}

	compFirst := v.Compare(vers[0])
	if compFirst == ComesBefore {
		return append([]SemanticVersion{v}, vers...)
	}

	compLast := v.Compare(vers[l-1])
	if compLast == ComesAfter {
		return append(vers, v)
	}

	ret := []SemanticVersion{}
	for i := 0; i < l-1; i++ {
		v1 := vers[i]
		ret = append(ret, v1)
		v2 := vers[i+1]
		if v1.Compare(v) == ComesBefore && v.Compare(v2) == ComesBefore {
			ret = append(ret, v)
		}
	}
	ret = append(ret, vers[l-1])
	return ret
}

func (index LibraryIndex) Versions(lib LibraryName) []SemanticVersion {
	ret := []SemanticVersion{}
	for _, v := range index.libraries {
		if v.library.name == lib {
			ret = addVersion(ret, v.library.ver)
		}
	}
	return ret
}
