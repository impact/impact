package utils

import "encoding/json"
import "io/ioutil"

type VersionString string;
type LibraryName string;

type Dependency struct {
	Name LibraryName `json:"name"`
	Version VersionString `json:"version"`
}

type Version struct {
	Version VersionString `json:"version"`
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
	Tarball string `json:"tarball_url"`
	Zipball string `json:"zipball_url"`
	Path string `json:"path"`
	Dependencies []Dependency `json:"dependencies"`
	Sha string `json:"sha"`
};

func (v1 Version) Equals(v2 Version) bool {
	return v1.Version==v2.Version;
}

type Library struct {
	Homepage string `json:"homepage"`
	Description string `json:"description"`
	Versions map[VersionString]Version `json:"versions"`
};

type Libraries map[LibraryName]Version;

type Index map[LibraryName]Library;

func (index *Index) ReadIndex(name string) error {
	file, err := ioutil.ReadFile(name)
	str := string(file)
	if (err==nil) {
		err = json.Unmarshal([]byte(str), index)
	}
	return err;
}

func (index *Index) Find(name LibraryName, version VersionString) (v Version, err error) {
	me, ok := (*index)[name];
	if (!ok) { err = MissingLibraryError{Name: name}; return; }
	v, ok = me.Versions[version];
	if (!ok) { err = MissingVersionError{Name: name, Version: version}; return; }
	return
}

func (index *Index) Dependencies(name LibraryName,
	                             version VersionString) (libs Libraries, err error) {
	// Find the specified library and find its version information
	// Then loop over all its dependencies and find call this function
	//   recursively for those.
	// Then collect all the individual libraries and check for version conflicts
	myver, err := index.Find(name, version);
	if (err!=nil) { return; }

	libs = Libraries{name: myver};

	for _, dep := range myver.Dependencies {
		deps, lerr := index.Dependencies(dep.Name, dep.Version);
		if (lerr!=nil) { return; }
		for dn, dv := range deps {
			ev, exists := libs[dn];
			if (exists) {
				if (dv.Equals(ev)) {
					err = VersionMismatchError{Name: dn, Existing: ev, Additional: dv};
					return;
				}
			} else {
				libs[dn] = dv;
			}
		}
	}
	return;
}
