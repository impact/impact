package index

import (
	"log"

	"github.com/wsxiaoys/terminal/color"

	"github.com/xogeny/impact/graph"
	"github.com/xogeny/impact/parsing"
)

func (ind *Index) BuildGraph(verbose bool) (graph.Resolver, error) {
	var resolver graph.Resolver = graph.NewLibraryGraph()

	// First, we collect all known libraries (these are essentially
	// all the potential nodes in the graph)
	for _, lib := range ind.Libraries {
		name := graph.LibraryName(lib.Name)
		for _, version := range lib.Versions {
			sver := version.Version
			resolver.AddLibrary(name, sver)
		}
	}

	// Now, we include dependencies between libraries (this are
	// the edges in the graph)
	for _, lib := range ind.Libraries {
		name := graph.LibraryName(lib.Name)
		for _, version := range lib.Versions {
			sver := version.Version
			for _, dependency := range version.Dependencies {
				dname := graph.LibraryName(dependency.Name)
				dver := dependency.Version

				dsver, err := parsing.NormalizeVersion(dver)
				if err != nil {
					log.Printf("Error parsing version %s: %v", dver, err)
					continue
				}

				// We ignore any errors (and hence, any dependency)
				err = resolver.AddDependency(name, sver, dname, dsver)
				if err != nil {
					if verbose {
						color.Println("@{r}Invalid dependency between:")
						color.Printf("  @{!r}%s %s @{r} and unknown library \n", name, sver)
						color.Printf("  @{!r}%s %s @{r}\n", dname, dsver)
					}
					continue
				}
			}
		}
	}

	return resolver, nil
}
