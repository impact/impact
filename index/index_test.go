package index

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func Test_Creation(t *testing.T) {
	Convey("Testing index creation", t, func(c C) {
		dep := Dependency{Name: "Foo", Version: "1.0.0"}
		version := Version{
			Version:      "0.1.0",
			Major:        0,
			Minor:        1,
			Patch:        0,
			Tarball:      "http://modelica.org/",
			Zipball:      "http://modelica.org/",
			Path:         "./ThisLibrary",
			Dependencies: []Dependency{dep},
			Sha:          "abcdefg",
		}
		lib := Library{
			Homepage:    "http://mylib.modelica.org",
			Description: "A dummy library",
			Versions:    map[VersionString]Version{"0.1.0": version},
		}

		index := map[string]Library{"Dummy": lib}

		_, ok := index["Dummy"]
		Equals(c, ok, true)
	})
}

func Test_UnmarshallDependency(t *testing.T) {
	Convey("Testing unmarshaling of dependencies", t, func(c C) {
		var ds = `{
                  "version": "3.2",
                  "name": "Modelica"
              }`
		sample := []byte(ds)
		dep := Dependency{}
		err := json.Unmarshal(sample, &dep)
		NoError(c, err)
		Equals(c, dep.Name, "Modelica")
		Equals(c, dep.Version, "3.2")
	})
}

func Test_UnmarshallVersion(t *testing.T) {
	Convey("Testing unmarshaling of versions", t, func(c C) {
		var ds = `{
                "major": 1,
                "tarball_url": "https://github.com/modelica-3rdparty/Physiolibrary/archive/v1.1.0.tar.gz",
                "patch": 0,
                "sha": "3075b23c214b65a510eb58654464f54507901378",
                "version": "1.1.0",
                "dependencies": [
                    {
                        "version": "3.2",
                        "name": "Modelica"
                    }
                ],
                "path": "Physiolibrary 1.1.0",
                "zipball_url": "https://github.com/modelica-3rdparty/Physiolibrary/archive/v1.1.0.zip",
                "minor": 1
            }`
		sample := []byte(ds)
		dep := Version{}
		err := json.Unmarshal(sample, &dep)
		NoError(c, err)
		Equals(c, dep.Major, 1)
		Equals(c, dep.Minor, 1)
		Equals(c, dep.Patch, 0)
		Equals(c, dep.Sha, "3075b23c214b65a510eb58654464f54507901378")
		Equals(c, dep.Version, "1.1.0")
		Equals(c, dep.Path, "Physiolibrary 1.1.0")
	})
}

func Test_UnmarshallLibrary(t *testing.T) {
	Convey("Testing unmarshaling of libraries", t, func(c C) {
		var ds = `{
                  "homepage": "http://www.modelica.org",
                  "description": "A dummy library",
                  "versions": {}
              }`
		sample := []byte(ds)
		dep := Library{}
		err := json.Unmarshal(sample, &dep)
		NoError(c, err)
		Equals(c, dep.Homepage, "http://www.modelica.org")
		Equals(c, dep.Description, "A dummy library")
	})
}

func Test_UnmarshallIndex(t *testing.T) {
	Convey("Testing unmarshaling of index", t, func(c C) {
		var ds = `{
    "SPICELib": {
        "homepage": "https://github.com/modelica-3rdparty/SPICELib",
        "description": "Free library with some of the modeling and analysis capabilities of the electric circuit simulator PSPICE.",
        "versions": {
            "1.1.0": {
                "major": 1,
                "tarball_url": "https://github.com/modelica-3rdparty/SPICELib/archive/v1.1.tar.gz",
                "patch": 0,
                "sha": "3d5738757b30192182b0b7caf46248c477d83e98",
                "version": "1.1.0",
                "dependencies": [],
                "path": "SPICELib 1.1",
                "zipball_url": "https://github.com/modelica-3rdparty/SPICELib/archive/v1.1.zip",
                "minor": 1
            },
            "1.1": {
                "major": 1,
                "tarball_url": "https://github.com/modelica-3rdparty/SPICELib/archive/v1.1.tar.gz",
                "patch": 0,
                "sha": "3d5738757b30192182b0b7caf46248c477d83e98",
                "version": "1.1.0",
                "dependencies": [],
                "path": "SPICELib 1.1",
                "zipball_url": "https://github.com/modelica-3rdparty/SPICELib/archive/v1.1.zip",
                "minor": 1
            }
        }
    },
    "FaultTriggering": {
        "homepage": "https://github.com/DLR-SR/FaultTriggering",
        "description": "Library for fault modelling in Modelica",
        "versions": {
            "0.6.2": {
                "major": 0,
                "tarball_url": "https://github.com/modelica-3rdparty/FaultTriggering/archive/v0.6.2.tar.gz",
                "patch": 2,
                "sha": "0a180687231d36540e1695523d3de6bbe10b28c5",
                "version": "0.6.2",
                "dependencies": [
                    {
                        "version": "3.2.1",
                        "name": "Modelica"
                    },
                    {
                        "version": "1.1.1",
                        "name": "ModelManagement"
                    }
                ],
                "path": "FaultTriggering",
                "zipball_url": "https://github.com/modelica-3rdparty/FaultTriggering/archive/v0.6.2.zip",
                "minor": 6
            },
            "0.5.0": {
                "major": 0,
                "tarball_url": "https://github.com/modelica-3rdparty/FaultTriggering/archive/v0.5.0.tar.gz",
                "patch": 0,
                "sha": "ad0a7ca17684753ceb3a6d5d2afd9988dc74912b",
                "version": "0.5.0",
                "dependencies": [
                    {
                        "version": "3.2.1",
                        "name": "Modelica"
                    },
                    {
                        "version": "1.1.1",
                        "name": "ModelManagement"
                    }
                ],
                "path": "FaultTriggering 0.5.0",
                "zipball_url": "https://github.com/modelica-3rdparty/FaultTriggering/archive/v0.5.0.zip",
                "minor": 5
            }
        }
    }
              }`
		sample := []byte(ds)
		dep := Index{}
		err := json.Unmarshal(sample, &dep)
		NoError(c, err)
	})
}

func Test_ReadFile(t *testing.T) {
	Convey("Testing reading of index files", t, func(c C) {
		index := Index{}
		err := index.BuildIndexFromFile("sample.json")
		NoError(c, err)
		_, ok := index["Physiolibrary"]
		Equals(c, ok, true)
	})
}

func contains(c C, libs Libraries, name LibraryName, ver VersionString) {
	c.Printf("Should contain library %s\n", name)
	v, ok := libs[name]
	IsTrue(c, ok)

	if v.Version != ver {
		c.Printf("Expected version %s of library %s but found %s\n",
			ver, name, v.Version)
	}
	Equals(c, v.Version, ver)
}

func Test_Dependencies(t *testing.T) {
	Convey("Testing dependencies", t, func(c C) {
		index := Index{}
		err := index.BuildIndexFromFile("sample.json")
		deps, err := index.Dependencies("MotorcycleLib", "1.0")
		NoError(c, err)
		for k, v := range deps {
			fmt.Println(k)
			fmt.Println(v.Version)
		}
		fmt.Println(deps)
		contains(c, deps, "MotorcycleLib", "1.0.0")
		contains(c, deps, "MultiBondLib", "1.3.0")
		contains(c, deps, "Modelica", "2.2.1")
		//contains(c, deps, "WheelsAndTires", "1.0.0")
	})
}
