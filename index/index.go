package index

import (
	"fmt"

	"encoding/json"
	"github.com/xogeny/impact/recorder"
)

type Index struct {
	Version   string     `json:"version"`
	Libraries []*Library `json:"libraries"`
}

func (i *Index) GetLibrary(owner string, name string) recorder.LibraryRecorder {
	for _, lib := range i.Libraries {
		if lib.Owner == owner && lib.Name == name {
			return lib
		}
	}
	lib := NewLibrary(owner, name)
	i.Libraries = append(i.Libraries, lib)
	return lib
}

func (i Index) JSON() (string, error) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewIndex() *Index {
	return &Index{
		Version:   "1.0.0",
		Libraries: []*Library{},
	}
}

var _ recorder.Recorder = (*Index)(nil)

func DownloadIndex() (Index, error) {
	/*
		    var master = "https://impact.modelica.org/impact_data.json"

			resp, err := http.Get(master)
			if err != nil {
				fmt.Println("Unable to locate index file at " + master)
				os.Exit(1)
			}
			defer resp.Body.Close()

			index := Index{}

			err = index.BuildIndex(resp.Body)
			if err != nil {
				fmt.Println("Error reading index: " + err.Error())
				os.Exit(2)
			}
	*/

	return Index{}, fmt.Errorf("Not implemented: Index.DownloadIndex")
}

/*
func (index *Index) BuildIndexFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return index.BuildIndex(file)
}

func (index *Index) BuildIndex(read io.Reader) error {
	data, err := ioutil.ReadAll(read)
	str := string(data)

	latest := Index{}
	if err == nil {
		err = json.Unmarshal([]byte(str), &latest)
	}
	for lib, v := range latest {
		v.Name = lib // Let the library know its name
		(*index)[lib] = v
	}
	return err
}

func (index *Index) Find(name LibraryName, version VersionString) (v VersionDetails, err error) {
	me, ok := (*index)[name]
	if !ok {
		err = MissingLibraryError{Name: name}
		return
	}
	v, ok = me.Versions[version]
	if !ok {
		err = MissingVersionError{Name: name, Version: version}
		return
	}
	return
}

// This function reports **only** a given libraries direct
// dependencies.  To resolve complete dependencies (including
// constraints due to other libraries, including dependencies,
// use the functionality in the graph package.
func (index *Index) Dependencies(name LibraryName,
	version VersionString) ([]Dependency, error) {

	myver, err := index.Find(name, version)
	if err != nil {
		return []Dependency{}, err
	}

	return myver.Dependencies, nil
}
*/
