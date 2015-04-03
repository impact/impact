package index

type Version struct {
	Version      VersionString `json:"version"`
	Major        int           `json:"major"`
	Minor        int           `json:"minor"`
	Patch        int           `json:"patch"`
	Tarball      string        `json:"tarball_url"`
	Zipball      string        `json:"zipball_url"`
	Path         string        `json:"path"`
	Dependencies []Dependency  `json:"dependencies"`
	Sha          string        `json:"sha"`
}

func (v1 Version) Equals(v2 Version) bool {
	// Perhaps this should do something about "-blah" and "+blah" data in the
	// version string?!?
	return v1.Version == v2.Version
}

func (v1 Version) GreaterThan(v2 Version) bool {
	return v1.Major > v2.Major ||
		(v1.Major == v2.Major && v1.Minor > v2.Minor) ||
		(v1.Major == v2.Major && v1.Minor == v2.Minor && v1.Patch > v2.Patch)
}
