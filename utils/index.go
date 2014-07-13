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

type MissingLibraryError struct {
	Name LibraryName
}

func (e MissingLibraryError) Error() string {
	return "No library named '"+string(e.Name)+"' found";
}

type MissingVersionError struct {
	Name LibraryName
	Version VersionString
}

func (e MissingVersionError) Error() string {
	return "No version '"+string(e.Version)+"' of library named '"+string(e.Name)+"' found";
}

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
	me, ok := (*index)[name];
	if (!ok) { err = MissingLibraryError{Name: name}; }
	myver, ok := me.Versions[version];
	if (!ok) { err = MissingVersionError{Name: name, Version: version}; }

	libs  = Libraries{name: myver};
	return;
}
