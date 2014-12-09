package deps

import "fmt"
import "log"
import "github.com/blang/semver"

/*
 * Create a special type to specifically represent library names.  This just
 * helps make the API clearer.
 */
type LibraryName string

/*
 * This type represents a specific configuration of libraries.  This is used to represent
 * the resolution of dependencies.
 */
type Configuration map[LibraryName]*semver.Version

func (conf Configuration) Clone() Configuration {
	clone := Configuration{}
	for k, v := range conf {
		clone[k] = v
	}
	return clone
}

/*
 * This type represents remaining possible values for a given library
 */
type Available map[LibraryName]*VersionList

func (a Available) Clone() Available {
	clone := Available{}
	for k, v := range a {
		clone[k] = v
	}
	return clone
}

func (a Available) Refine(subset Available) Available {
	ret := Available{}

	for k, v := range a {
		v2, exists := subset[k]
		if !exists {
			ret[k] = v.Clone()
		} else {
			ret[k] = (*v).Intersection(*v2)
		}
		ret[k].ReverseSort()
	}
	return ret
}

func (a Available) Empty() []LibraryName {
	ret := []LibraryName{}
	for k, v := range a {
		if v.Len() == 0 {
			ret = append(ret, k)
		}
	}
	return ret
}

/*
 * This struct is not exported.  It is used to represent a unique library
 * (i.e., name + version)
 */
type uniqueLibrary struct {
	name LibraryName
	ver  *semver.Version
}

func (l1 uniqueLibrary) Equals(lib LibraryName, ver *semver.Version) bool {
	return l1.name == lib && l1.ver.Compare(ver) == 0
}

/*
 * This is an edge in our (directed) dependency graph.  It indicates that `library`
 * depends on `dependsOn`.  Each is represented as a unique library (i.e., name + version)
 */
type dependency struct {
	library   uniqueLibrary
	dependsOn uniqueLibrary
}

/*
 * A library index is simply a list of dependencies (edges)
 */
type LibraryIndex struct {
	dependencies []dependency
	libraries    []uniqueLibrary
}

/*
 * This function creates a new LibraryIndex object.
 */
func MakeLibraryIndex() LibraryIndex {
	return LibraryIndex{
		libraries:    []uniqueLibrary{},
		dependencies: []dependency{},
	}
}

func (index LibraryIndex) Contains(lib LibraryName, libver *semver.Version) bool {
	for _, l := range index.libraries {
		if l.Equals(lib, libver) {
			return true
		}
	}
	return false
}

func (index *LibraryIndex) AddLibrary(lib LibraryName, libver *semver.Version) {
	for _, l := range index.libraries {
		if l.Equals(lib, libver) {
			return
		}
	}
	index.libraries = append(index.libraries, uniqueLibrary{name: lib, ver: libver})
}

/*
 * Method to add a new dependency to a library index
 */
func (index *LibraryIndex) AddDependency(lib LibraryName, libver *semver.Version,
	deplib LibraryName, depver *semver.Version) error {

	lfound := false
	dfound := false

	for _, l := range index.libraries {
		if l.Equals(lib, libver) {
			lfound = true
		}
		if l.Equals(deplib, depver) {
			dfound = true
		}
	}

	if !lfound {
		return fmt.Errorf("Cannot add dependency for unknown library %s v%s",
			lib, libver.String())
	}
	if !dfound {
		return fmt.Errorf("Cannot add dependency for unknown library %s v%s",
			deplib, depver.String())
	}

	library := uniqueLibrary{name: lib, ver: libver}
	dependsOn := uniqueLibrary{name: deplib, ver: depver}
	dep := dependency{library: library, dependsOn: dependsOn}

	index.dependencies = append(index.dependencies, dep)
	return nil
}

/*
 * Builds a list of all versions of a given library known to the
 * index.  These are returned in sorted order (latest to earliest)
 */
func (index LibraryIndex) Versions(lib LibraryName) *VersionList {
	present := map[*semver.Version]bool{}

	for _, l := range index.libraries {
		if l.name == lib {
			present[l.ver] = true
		}
	}

	vl := NewVersionList()
	for v, _ := range present {
		vl.Add(v)
	}

	vl.ReverseSort()
	return vl
}

/*
 * This method
 */
func (index LibraryIndex) Dependencies(lib LibraryName, ver *semver.Version) Available {
	depvers := Available{}

	for _, dep := range index.dependencies {
		// Is this a dependency for the current library and version?
		if dep.library.name == lib && ver.Compare(dep.library.ver) == 0 {
			// If so, add it to the available set (if one exists)
			dver, found := depvers[dep.dependsOn.name]
			if !found {
				dver = NewVersionList()
				depvers[dep.dependsOn.name] = dver
			}
			dver.Add(dep.dependsOn.ver)
		}
	}
	return depvers
}

func (index LibraryIndex) findFirst(
	mapped Configuration, // Variables whose values have already been chosen
	verbose bool, // Whether to generate verbose output
	avail Available, // Constraints of possible values for remaining libraries
	rest ...LibraryName, // Libraries whose versions we still need to decide
) (Configuration, error) {
	if verbose {
		log.Printf("Call to findFirst...")
		log.Printf("  Mapped: %v", mapped)
		log.Printf("  Avail: %v", avail)
		log.Printf("  Rest: %v", rest)
	}

	// Nothing left to process...we are done!
	if len(rest) == 0 {
		if verbose {
			log.Printf("End of the line, returning %v", mapped)
		}
		return mapped, nil
	}

	// Consider the next library in the list
	lib := rest[0]
	rest = rest[1:]

	if verbose {
		log.Printf("  -> Lib = %v", lib)
		log.Printf("  -> Rest = %v", rest)
	}

	// Determine remaining possible values for this library...
	vers := avail[lib]

	// Loop over each possible version of the chosen library
	for _, ver := range *vers {
		if verbose {
			log.Printf("  Considering version %v of %s", ver, lib)
		}

		/* Create our own local copy of the configuration so we don't mutate 'mapped' */
		config := mapped.Clone()
		// A list of any new libraries to introduce to the search
		newlibs := []LibraryName{}

		// Find out all the libraries that this particular library+version depend on
		depvers := index.Dependencies(lib, ver)

		if verbose {
			log.Printf("Dependencies of %s:%s -> %s", lib, ver.String(), depvers)
		}

		// Have any of this libraries dependencies already been chosen?
		for d, vl := range depvers {
			choice, chosen := mapped[d]
			if chosen {
				// If our choice is not among the set that this library depends on,
				// we are done.
				if !vl.Contains(choice) {
					return nil, fmt.Errorf("No compatible version of %s", d)
				}
				// Otherwise, the current choice is compatible
			}
		}

		// Ignore any previous mapped libraries (we just checked to make sure
		// we were compatible with those in the previous few lines of code so
		// we can safely ignore them)
		for l, _ := range mapped {
			delete(depvers, l)
		}

		if verbose {
			log.Printf("Dependencies after mapping -> %s", depvers)
		}

		newavail := avail.Clone()

		// Add any new dependencies?  (Check to see if we were already planning on
		// incuding them, if not add them)
		for n1, _ := range depvers {
			found := false
			for _, n2 := range rest {
				if n1 == n2 {
					found = true
				}
			}
			if !found {
				newlibs = append(newlibs, n1)
				_, exists := newavail[n1]
				if exists {
					panic(fmt.Errorf("Should not happen"))
				}
				newavail[n1] = index.Versions(n1)
			}
		}

		// Take the intersection of the previously available versions with
		// the dependent versions
		newavail = newavail.Refine(depvers)

		if verbose {
			log.Printf("Intersection between %s and %s -> %s", avail, depvers, newavail)
		}

		// Make sure the current library is removed from this list
		delete(newavail, lib)

		if verbose {
			log.Printf("Intersection - %s -> %s", lib, newavail)
		}

		// Are any of the available value sets empty?  If so, return an error
		empty := newavail.Empty()
		if len(empty) > 0 {
			return nil, fmt.Errorf("No compatible versions of: %v", empty)
		}

		// Specify the current library and version choice
		config[lib] = ver

		// Recurse to solve remaining variables
		newlibs = append(newlibs, rest...)
		sol, err := index.findFirst(config, verbose, newavail, newlibs...)
		if err == nil {
			return sol, err
		}
	}
	return nil, fmt.Errorf("No compatible versions of %s found", lib)
}

func (index LibraryIndex) Resolve(libraries ...LibraryName) (config Configuration, err error) {
	// Initially, all possible values for the libraries are possible
	ret := Available{}
	for _, lib := range libraries {
		ret[lib] = index.Versions(lib)
	}

	// Now search for a consistent set...
	return index.findFirst(config, false, ret, libraries...)
}
