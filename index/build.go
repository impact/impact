package index

import (
	"log"

	"github.com/xogeny/impact/graph"
	"github.com/xogeny/impact/parsing"
)

func (ind *Index) BuildGraph(verbose bool) (graph.Resolver, error) {
	var resolver graph.Resolver = graph.NewLibraryGraph()

	for _, lib := range ind.Libraries {
		name := graph.LibraryName(lib.Name)
		for _, version := range lib.Versions {
			sver := version.Version
			resolver.AddLibrary(name, sver)
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
						log.Printf("Invalid dependency of %s:%v on %s:%v: %v",
							name, sver, dname, dsver, err)
					}
					continue
				}
			}
		}
	}

	return resolver, nil
}
