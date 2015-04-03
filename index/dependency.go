package index

type Dependency struct {
	Name    LibraryName   `json:"name"`
	Version VersionString `json:"version"`
}
