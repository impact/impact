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

type Index map[string]Library;

func ReadIndex(name string, index *Index) error {
	file, err := ioutil.ReadFile(name)
	str := string(file)
	if (err==nil) {
		err = json.Unmarshal([]byte(str), index)
	}
	return err;
}

func (index *Index) Dependencies(name LibraryName, version VersionString) Libraries {
	return Libraries{};
}
