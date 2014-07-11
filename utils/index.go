package utils

import "encoding/json"
import "io/ioutil"

type Dependency struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

type Version struct {
	Version string `json:"version"`
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
	Tarball string `json:"tarball_url"`
	Zipball string `json:"zipball_url"`
	Path string `json:"path"`
	Dependencies []Dependency `json:"dependencies"`
	Sha string `json:"sha"`
};

type VersionString string;
type LibraryName string;

type Library struct {
	Homepage string `json:"homepage"`
	Description string `json:"description"`
	Versions map[VersionString]Version `json:"versions"`
};

type Libraries map[LibraryName]Version;

type Index map[LibraryName]Library;

func ReadIndex(name string, index *Index) error {
	file, err := ioutil.ReadFile(name)
	str := string(file)
	if (err==nil) {
		err = json.Unmarshal([]byte(str), index)
	}
	return err;
}

func (index *Index) Dependencies(name LibraryName,
	                             version VersionString) (libs Libraries, err error) {
	// Find the specified library and find its version information
	// Then loop over all its dependencies and find call this function
	//   recursively for those.
	// Then collect all the individual libraries and check for version conflicts
	me, _ := (*index)[name];
	// TODO: Check that it is ok
	myver, _ := me.Versions[version];
	// TODO: Check that it is ok

	libs  = Libraries{name: myver};
	return;
}
