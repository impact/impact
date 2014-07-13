package main

import "testing"
import "xogeny/gimpact/utils"
import "encoding/json"
import "github.com/stretchr/testify/assert"
import "fmt"

func Test_Creation(t* testing.T) {
	dep := utils.Dependency{Name: "Foo", Version: "1.0.0"};
	version := utils.Version{
		Version: "0.1.0",
		Major: 0,
		Minor: 1,
		Patch: 0,
		Tarball: "http://modelica.org/",
		Zipball: "http://modelica.org/",
		Path: "./ThisLibrary",
		Dependencies: []utils.Dependency{dep},
		Sha: "abcdefg",
	};
	lib := utils.Library{
		Homepage: "http://mylib.modelica.org",
		Description: "A dummy library",
		Versions: map[utils.VersionString]utils.Version{"0.1.0": version},
	};

	index := map[string]utils.Library{"Dummy": lib}

	_, ok := index["Dummy"]
	assert.Equal(t, ok, true, "Library not in index");
}

func Test_UnmarshallDependency(t* testing.T) {
	var ds = `{
                  "version": "3.2", 
                  "name": "Modelica"
              }`
	sample := []byte(ds);
	dep := utils.Dependency{};
	err := json.Unmarshal(sample, &dep);
	assert.NoError(t, err);
	assert.Equal(t, dep.Name, "Modelica", "name mismatch");
	assert.Equal(t, dep.Version, "3.2", "version");
}

func Test_UnmarshallVersion(t* testing.T) {
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
            }`;
	sample := []byte(ds);
	dep := utils.Version{};
	err := json.Unmarshal(sample, &dep);
	assert.NoError(t, err);
	assert.Equal(t, dep.Major, 1, "major mismatch");
	assert.Equal(t, dep.Minor, 1, "minor mismatch");
	assert.Equal(t, dep.Patch, 0, "patch mismatch");
	assert.Equal(t, dep.Sha, "3075b23c214b65a510eb58654464f54507901378", "sha mismatch");
	assert.Equal(t, dep.Version, "1.1.0", "version mismatch");
	assert.Equal(t, dep.Path, "Physiolibrary 1.1.0", "path mismatch");
}

func Test_UnmarshallLibrary(t* testing.T) {
	var ds = `{
                  "homepage": "http://www.modelica.org",
                  "description": "A dummy library",
                  "versions": {}
              }`
	sample := []byte(ds);
	dep := utils.Library{};
	err := json.Unmarshal(sample, &dep);
	assert.NoError(t, err);
	assert.Equal(t, dep.Homepage, "http://www.modelica.org", "homepage mismatch");
	assert.Equal(t, dep.Description, "A dummy library", "description mismatch");
}

func Test_UnmarshallIndex(t* testing.T) {
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
	sample := []byte(ds);
	dep := utils.Index{};
	err := json.Unmarshal(sample, &dep);
	assert.NoError(t, err);
}

func Test_ReadFile(t* testing.T) {
	index := utils.Index{};
	err := index.BuildIndexFromFile("sample.json");
	assert.NoError(t, err);
	_, ok := index["Physiolibrary"];
	assert.Equal(t, ok, true, "Couldn't find Physiolibrary");
}

func contains(t* testing.T, libs utils.Libraries, name utils.LibraryName, ver utils.VersionString) {
	v, ok := libs[name];
	if (!ok) { t.Fatal("Version "+string(ver)+" of library "+string(name)+" not found"); }
	if (v.Version!=ver) {
		t.Fatal("Expected version "+string(ver)+" of library "+string(name)+
			" but found "+string(v.Version));
	}
}

func Test_Dependencies(t* testing.T) {
	index := utils.Index{};
	err := index.BuildIndexFromFile("sample.json");
	deps, err := index.Dependencies("MotorcycleLib", "1.0");
	assert.NoError(t, err);
	for k, v := range deps {
		fmt.Println(k);
		fmt.Println(v.Version);
	}
	fmt.Println(deps);
	contains(t, deps, "MotorcycleLib", "1.0.0");
	contains(t, deps, "MultiBondLib", "1.3.0");
	contains(t, deps, "Modelica", "2.2.1");
	contains(t, deps, "WheelsAndTires", "1.0.0");
}
