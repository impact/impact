package main

import "testing"
import "xogeny/gimpact/utils"
import "encoding/json"
import "github.com/stretchr/testify/assert"

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
	assert.Nil(t, err, "Unmarshal failed");
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
	assert.Nil(t, err, "Unmarshal failed");
	if (dep.Major!=1) { t.Fatal("major mismatch"); }
	if (dep.Minor!=1) { t.Fatal("minor mismatch"); }
	if (dep.Patch!=0) { t.Fatal("patch mismatch"); }
	if (dep.Sha!="3075b23c214b65a510eb58654464f54507901378") { t.Fatal("sha mismatch"); }
	if (dep.Version!="1.1.0") { t.Fatal("version mismatch"); }
	if (dep.Path!="Physiolibrary 1.1.0") { t.Fatal("version mismatch"); }
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
	if (err!=nil) { t.Fatal("Unmarshal failed"); }
	if (dep.Homepage!="http://www.modelica.org") { t.Fatal("homepage mismatch"); }
	if (dep.Description!="A dummy library") { t.Fatal("description mismatch"); }
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
	if (err!=nil) { t.Fatal("Unmarshal failed"); }
}

func Test_ReadFile(t* testing.T) {
	index := utils.Index{};
	err := utils.ReadIndex("sample.json", &index);
	if (err!=nil) { t.Fatal("Error reading file: "+err.Error()); }
	_, ok := index["Physiolibrary"];
	if (!ok) { t.Fatal("Couldn't find Physiolibrary"); }
}
