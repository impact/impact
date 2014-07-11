package impact

import "testing"

func Test_Creation(t* testing.T) {
	dep := Dependency{Name: "Foo", Version: "1.0.0"};
	version := Version{
		Version: "0.1.0",
		Major: 0,
		Minor: 1,
		Patch: 0,
		Tarball: "http://modelica.org/",
		Zipball: "http://modelica.org/",
		Path: "./ThisLibrary",
		Dependencies: []Dependency{dep},
		Sha: "abcdefg",
	};
	lib := Library{
		Homepage: "http://mylib.modelica.org",
		Description: "A dummy library",
		Versions: map[string]Version{"0.1.0": version},
	};

	index := map[string]Library{"Dummy": lib}

	_, ok := index["Dummy"]
	if (!ok) { t.Fatal("Library not in index"); }
}
