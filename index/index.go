package index

import (
	"encoding/json"
	"fmt"

	"github.com/blang/semver"

	"github.com/xogeny/impact/recorder"
)

type Index struct {
	Version   string     `json:"version"`
	Libraries []*Library `json:"libraries"`
}

func (i Index) Find(name string, version semver.Version) (VersionDetails, error) {
	for _, lib := range i.Libraries {
		if lib.Name == name {
			for _, details := range lib.Versions {
				if details.Version.EQ(version) {
					return *details, nil
				}
			}
		}
	}
	return VersionDetails{},
		fmt.Errorf("Couldn't find version %v of library %s",
			version, name)
}

func (i *Index) GetLibrary(name string, uri string, owner_uri string) recorder.LibraryRecorder {
	for _, lib := range i.Libraries {
		if lib.OwnerURI == owner_uri && lib.Name == name {
			return lib
		}
	}
	lib := NewLibrary(name, uri, owner_uri)
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
