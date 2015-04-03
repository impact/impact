package graph

import (
	"fmt"
	"log"

	"github.com/blang/semver"
)

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

func (d dependency) String() string {
	return fmt.Sprintf("%s:%s -> %s:%s", d.library.name,
		d.library.ver.String(), d.dependsOn.name,
		d.dependsOn.ver.String())
}

func (d1 dependency) Equals(d2 dependency) bool {
	return d1.library.name == d2.library.name &&
		d1.library.ver.Compare(d2.library.ver) == 0 &&
		d1.dependsOn.name == d2.dependsOn.name &&
		d1.dependsOn.ver.Compare(d2.dependsOn.ver) == 0
}

/*
 * A library index is simply a list of dependencies (edges)
 */
type LibraryGraph struct {
	dependencies []dependency
	libraries    []uniqueLibrary
	verbose      bool
}

func (index LibraryGraph) Exists(d dependency) bool {
	for _, c := range index.dependencies {
		if c.Equals(d) {
			return true
		}
	}
	return false
}

/*
 * This function creates a new LibraryGraph object.
 */
func NewLibraryGraph() *LibraryGraph {
	return &LibraryGraph{
		libraries:    []uniqueLibrary{},
		dependencies: []dependency{},
		verbose:      false,
	}
}

/*
 * This function sets debugging output
 */
func (index *LibraryGraph) Verbose(v bool) {
	index.verbose = v
}

func (index LibraryGraph) Contains(lib LibraryName, libver *semver.Version) bool {
	for _, l := range index.libraries {
		if l.Equals(lib, libver) {
			return true
		}
	}
	return false
}

func (index *LibraryGraph) AddLibrary(lib LibraryName, libver *semver.Version) {
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
func (index *LibraryGraph) AddDependency(lib LibraryName, libver *semver.Version,
	deplib LibraryName, depver *semver.Version) error {

	/* Flag indicating whether we have found the `lib` library */
	lfound := false
	/* Flag indicating whether we have found the `deplib` library */
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

	// Avoid duplicates
	// TODO: Get rid of this when quality of index data is such that we don't
	// need it...
	if !index.Exists(dep) {
		index.dependencies = append(index.dependencies, dep)
	}
	return nil
}

/*
 * Builds a list of all versions of a given library known to the
 * index.  These are returned in sorted order (latest to earliest)
 */
func (index LibraryGraph) Versions(lib LibraryName) *VersionList {
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
 * This method returns a map where the keys are libraries that the named
 * library+version depends on and the values are the versions that are
 * compatible with the named library+version.
 */
func (index LibraryGraph) Dependencies(lib LibraryName, ver *semver.Version) Possible {
	depvers := Possible{}

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

func (index LibraryGraph) findFirst(
	mapped Configuration, // Variables whose values have already been chosen
	verbose bool, // Whether to generate verbose output
	avail Possible, // Constraints of possible values for remaining libraries
	rest ...LibraryName, // Libraries whose versions we still need to decide
) (Configuration, error) {
	if verbose {
		log.Printf("Call to findFirst...")
		log.Printf("  Resolved: %v", mapped)
		log.Printf("  Avail: %v", avail)
		log.Printf("  Yet to be Resolved: %v", rest)
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
		log.Printf("  -> Library to Resolve: %v", lib)
		log.Printf("  -> Remaining: %v", rest)
	}

	// Determine remaining possible values for this library...
	vers := avail[lib]
	if vers == nil {
		log.Printf("Error, unable to find library %s (mapped: %v, rest: %v)",
			lib, mapped, rest)
		return nil, fmt.Errorf("Internal error handling library %s", lib)
	}

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

		// Check to see if any of these dependencies have already been chosen
		// (Note: I tried to make this process a separate function and CPU
		//  time doubled because of something to do with semaphores...odd)
		for d, vl := range depvers {
			choice, chosen := mapped[d]
			if chosen {
				// If our choice is not among the set that this library depends on,
				// there is no possible solution and we are done here.
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
		for n1, v1 := range depvers {
			if n1 == lib {
				if v1.Contains(ver) {
					// We can skip n1 because it is a *consistent* circular dependency
					// that will automatically be met
					continue
				} else {
					return nil,
						fmt.Errorf("Inconsistent circular dependency for library %v", lib)
				}
			}
			found := false
			for _, n2 := range rest {
				if n1 == n2 {
					found = true
				}
			}
			if !found {
				// If we get here, then we have a key in depvers that is not present
				// in the 'rest' list.  So we have a new variable to solve for...
				newlibs = append(newlibs, n1)
				// Our initial set of possible values for a new library is all versions
				// in the index.
				newavail[n1] = index.Versions(n1)
			}
		}

		// Refine the set of possible values for variables based on what
		// is permitted by our dependencies
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
		// since it means that there will be no solution for those variables
		empty := newavail.Empty()
		if len(empty) > 0 {
			return nil, fmt.Errorf("No compatible versions of: %v", empty)
		}

		// Specify the current library and version choice
		config[lib] = ver

		// TODO: Check newlibs to see if the current library is in there
		// and make sure that the choice above didn't conflict

		// Prepend any new libraries to the front of the list of libraries
		// we still need to solve for (this makes things depth first)
		newlibs = append(newlibs, rest...)

		// Recurse to solve remaining variables
		sol, err := index.findFirst(config, verbose, newavail, newlibs...)
		if err == nil {
			return sol, err
		}
	}
	return nil, fmt.Errorf("No compatible versions of %s found", lib)
}

func (index LibraryGraph) Resolve(libraries ...LibraryName) (config Configuration, err error) {
	// Initially, all possible values for the libraries are possible
	ret := Possible{}
	for _, lib := range libraries {
		ret[lib] = index.Versions(lib)
	}

	// Now search for a consistent set...
	return index.findFirst(config, index.verbose, ret, libraries...)
}
