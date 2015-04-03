package graph

import (
	"fmt"
	"strings"

	"testing"

	"github.com/blang/semver"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

/*
 * Helper function to turn a string (in the form 'LibName:Version' into a
 * library name and semantic version.
 */
func parse(libs string) (LibraryName, semver.Version) {
	parts := strings.Split(libs, ":")
	if len(parts) != 2 {
		panic(fmt.Errorf("Invalid library spec: %s", libs))
	}
	v, err := semver.New(parts[1])
	if err != nil {
		panic(err)
	}
	return LibraryName(parts[0]), v
}

/*
 * Declare a dependency and add libraries, all in one shot.
 */
func deps(index Resolver, libs string, deps ...string) error {
	lib, libver := parse(libs)
	index.AddLibrary(lib, libver)
	for _, ds := range deps {
		dep, depver := parse(ds)
		if !index.Contains(dep, depver) {
			index.AddLibrary(dep, depver)
		}
		err := index.AddDependency(lib, libver, dep, depver)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 * Test whether the resolution matches expected values.
 */
func testConfig(config Configuration, c C, vers ...string) {
	for _, v := range vers {
		lib, libver := parse(v)
		cver, exists := config[lib]
		IsTrue(c, exists)
		if exists {
			Equals(c, 0, cver.Compare(libver))
		}
	}
}

/* Simple Case: Root 1.0.0 depends on A 1.0.0 */
func TestResolution3(t *testing.T) {
	Convey("Testing very basic dependency of one library on another", t, func(c C) {
		var index Resolver = NewLibraryGraph()
		err := deps(index, "Root:1.0.0", "A:1.0.0")
		NoError(c, err)
		config, err := index.Resolve("Root")
		NoError(c, err)

		testConfig(config, c, "Root:1.0.0")
	})
}

/* Circular Case:
 *   Root 1.0.0 -> A 1.0.0 AND
 *   A 1.0.0    -> Root 1.0.0
 */
func TestResolutionOfCircularDependency(t *testing.T) {
	Convey("Testing circular dependencies", t, func(c C) {
		var index Resolver = NewLibraryGraph()

		err := deps(index, "Root:1.0.0", "A:1.0.0")
		NoError(c, err)

		err = deps(index, "A:1.0.0", "Root:1.0.0")
		NoError(c, err)

		config, err := index.Resolve("Root")
		NoError(c, err)

		testConfig(config, c, "Root:1.0.0", "A:1.0.0")
	})
}

/* Unmet Circular
 *   Root 1.0.0 -> A 1.0.0 AND
 *   A 1.0.0    -> Root 1.0.1
 *
 *   Root 1.0.1 -> A 1.0.1 AND
 *   A 1.0.1    -> Root 1.0.0
 */
func TestResolutionOfUnmetCircularDependency(t *testing.T) {
	Convey("Testing unmet circular dependencies", t, func(c C) {
		var index Resolver = NewLibraryGraph()

		err := deps(index, "Root:1.0.0", "A:1.0.0")
		NoError(c, err)

		err = deps(index, "A:1.0.0", "Root:1.0.1")
		NoError(c, err)

		err = deps(index, "Root:1.0.1", "A:1.0.1")
		NoError(c, err)

		err = deps(index, "A:1.0.1", "Root:1.0.0")
		NoError(c, err)

		/* Should yield an error, since no configuration works */
		_, err = index.Resolve("Root")
		IsError(c, err)
	})
}

/*
 * A 1.0.0 -> B 1.0.1 | B 1.0.0
 * C 1.0.0 -> D 1.0.1 | D 1.0.0
 * E 1.0.0 -> F 1.0.1 | F 1.0.0
 */
func TestResolutionSimple1(t *testing.T) {
	Convey("Testing slightly complex dependencies", t, func(c C) {
		var index Resolver = NewLibraryGraph()

		err := deps(index, "A:1.0.0", "B:1.0.1")
		NoError(c, err)

		err = deps(index, "A:1.0.0", "B:1.0.0")
		NoError(c, err)

		err = deps(index, "C:1.0.0", "D:1.0.1")
		NoError(c, err)

		err = deps(index, "C:1.0.0", "D:1.0.0")
		NoError(c, err)

		err = deps(index, "E:1.0.0", "F:1.0.1")
		NoError(c, err)

		err = deps(index, "E:1.0.0", "F:1.0.0")
		NoError(c, err)

		config, err := index.Resolve("A", "C", "E")
		NoError(c, err)
		testConfig(config, c, "A:1.0.0", "B:1.0.1", "C:1.0.0", "D:1.0.1", "E:1.0.0", "F:1.0.1")
	})
}

/*
 * A 1.0.0 -> B 1.0.1 | B 1.0.0
 * B 1.0.1 -> D 1.0.0
 * C 1.0.0 -> D 1.0.1 | D 1.0.0
 * E 1.0.0 -> F 1.0.1 | F 1.0.0
 */
func TestResolutionSimple2(t *testing.T) {
	Convey("Testing more slightly complex dependencies", t, func(c C) {
		var index Resolver = NewLibraryGraph()

		err := deps(index, "A:1.0.0", "B:1.0.1")
		NoError(c, err)

		err = deps(index, "A:1.0.0", "B:1.0.0")
		NoError(c, err)

		err = deps(index, "B:1.0.1", "D:1.0.0")
		NoError(c, err)

		err = deps(index, "C:1.0.0", "D:1.0.1")
		NoError(c, err)

		err = deps(index, "E:1.0.0", "F:1.0.1")
		NoError(c, err)

		err = deps(index, "E:1.0.0", "F:1.0.0")
		NoError(c, err)

		config, err := index.Resolve("A", "C", "E")
		NoError(c, err)

		c.Printf("config = %v\n", config)
		testConfig(config, c, "A:1.0.0", "B:1.0.0", "C:1.0.0", "D:1.0.1", "E:1.0.0", "F:1.0.1")
	})
}

/*
 * A 1.0.0 -> B 1.0.0
 * A 1.0.0 -> B 1.0.1
 * A 1.0.1 -> B 1.0.0
 * A 1.0.1 -> B 1.0.1
 * B 1.0.0 -> C 1.0.0
 * B 1.0.0 -> C 1.0.1
 * B 1.0.1 -> C 1.0.0
 * B 1.0.1 -> C 1.0.1
 * C 1.0.0 -> D 1.0.0
 * C 1.0.0 -> D 1.0.1
 * C 1.0.1 -> D 1.0.0
 * C 1.0.1 -> D 1.0.1
 *
 * Z 1.0.0 -> A 1.0.0, B 1.0.0, C 1.0.0, D 1.0.0
 */
func TestWorstCaseScenario(t *testing.T) {
	Convey("Testing worst case scenario of exhaustive searching", t, func(c C) {
		var index Resolver = NewLibraryGraph()

		solve := []LibraryName{}
		zdeps := []string{}
		letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
		//letters := []string{"A", "B", "C", "D"}
		for i, l := range letters {
			if i < len(letters)-1 {
				err := deps(index,
					fmt.Sprintf("%s:1.0.0", l),
					fmt.Sprintf("%s:1.0.0", letters[i+1]))
				NoError(c, err)

				err = deps(index,
					fmt.Sprintf("%s:1.0.0", l),
					fmt.Sprintf("%s:1.0.1", letters[i+1]))
				NoError(c, err)

				err = deps(index,
					fmt.Sprintf("%s:1.0.1", l),
					fmt.Sprintf("%s:1.0.0", letters[i+1]))
				NoError(c, err)

				err = deps(index,
					fmt.Sprintf("%s:1.0.1", l),
					fmt.Sprintf("%s:1.0.1", letters[i+1]))
				NoError(c, err)
			}
			solve = append(solve, LibraryName(l))
			zdeps = append(zdeps, fmt.Sprintf("%s:1.0.0", l))
		}

		err := deps(index, "Z:1.0.0", zdeps...)
		NoError(c, err)

		solve = append(solve, "Z")

		/* Should yield an error, since no configuration works */
		c.Printf("Solve for  = %v\n", solve)
		config, err := index.Resolve(solve...)
		NoError(c, err)

		c.Printf("config = %v\n", config)
		testConfig(config, c, zdeps...)
	})
}

func TestSelfDependence(t *testing.T) {
	Convey("Testing self dependencies", t, func(c C) {
		index := NewLibraryGraph()

		err := deps(index, "A:1.0.1", "A:1.0.1")
		NoError(c, err)

		err = deps(index, "A:1.0.0", "A:1.0.0")
		NoError(c, err)

		config, err := index.Resolve("A")
		NoError(c, err)

		testConfig(config, c, "A:1.0.1")
	})
}

/*
 * This case tests the lower level (pedantic) API.
 */
func TestResolution1(t *testing.T) {
	Convey("Testing low-level functionality", t, func(c C) {
		index := NewLibraryGraph()
		root1, err := semver.Parse("1.0.0")
		NoError(c, err)

		a1, err := semver.Parse("1.0.0")
		NoError(c, err)

		index.AddLibrary("Root", root1)
		err = index.AddDependency("Root", root1, "A", a1)
		IsError(c, err)

		index.AddLibrary("A", a1)
		err = index.AddDependency("Root", root1, "A", a1)
		NoError(c, err)

		rootVers := index.Versions("Root")
		Resembles(c, *rootVers, VersionList{root1})

		config, err := index.Resolve("Root")
		NoError(c, err)
		NotNil(c, config["Root"])

		c.Printf("Configuration: %v\n", config)

		// Introduce a circular dependency.  Make sure nothing breaks
		err = index.AddDependency("A", a1, "Root", root1)
		NoError(c, err)

		rootVers = index.Versions("Root")
		IsTrue(c, rootVers.Contains(root1))

		Equals(c, 1, rootVers.Len())

		config, err = index.Resolve("Root")
		NoError(c, err)

		c.Printf("Configuration: %v\n", config)
	})
}

/*
 * This case also tests the lower level (pedantic) API.
 */
func TestResolution2(t *testing.T) {
	Convey("Testing more low-level functionality", t, func(c C) {
		index := NewLibraryGraph()
		root1, err := semver.Parse("1.0.0")
		NoError(c, err)

		a1, err := semver.Parse("1.0.0")
		NoError(c, err)

		index.AddLibrary("Root", root1)
		index.AddLibrary("A", a1)
		err = index.AddDependency("Root", root1, "A", a1)
		NoError(c, err)

		rootVers := index.Versions("Root")
		Resembles(c, *rootVers, VersionList{root1})

		// Introduce a circular dependency that makes resolution fail
		root2, err := semver.New("1.0.1")
		a2, err := semver.New("1.0.1")
		NoError(c, err)

		index.AddLibrary("Root", root2)
		index.AddLibrary("A", a2)
		err = index.AddDependency("Root", root2, "A", a2)
		err = index.AddDependency("A", a1, "Root", root2)
		err = index.AddDependency("A", a2, "Root", root1)
		NoError(c, err)

		rootVers = index.Versions("Root")
		IsTrue(c, rootVers.Contains(root1))
		IsTrue(c, rootVers.Contains(root2))
		Equals(c, 2, rootVers.Len())

		_, err = index.Resolve("Root")
		IsError(c, err)
	})
}
