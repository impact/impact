package utils

import "encoding/json"
import "io/ioutil"
import "strings"
import "io"
import "os"

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
	Name LibraryName
	Homepage string `json:"homepage"`
	Description string `json:"description"`
	Versions map[VersionString]Version `json:"versions"`
};

func (lib Library) Matches(term string) bool {
	var match = false;
	if (strings.Contains(string(lib.Name), term)) {
		match = true;
	}
	if (strings.Contains(string(lib.Description), term)) {
		match = true;
	}
	return match;
}

type Libraries map[LibraryName]Version;

type Index map[LibraryName]Library;

func (index *Index) BuildIndexFromFile(filename string) error {
	file, err := os.Open(filename);
	if (err!=nil) { return err; }
	defer file.Close();
	return index.BuildIndex(file);
}

func (index *Index) BuildIndex(read io.Reader) error {
	data, err := ioutil.ReadAll(read);
	str := string(data)

	latest := Index{};
	if (err==nil) {
		err = json.Unmarshal([]byte(str), &latest)
	}
	for lib, v := range latest {
		v.Name = lib; // Let the library know its name
		(*index)[lib] = v;
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
