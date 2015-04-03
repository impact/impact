package index

type MissingLibraryError struct {
	Name LibraryName
}

func (e MissingLibraryError) Error() string {
	return "No library named '" + string(e.Name) + "' found"
}

type MissingVersionError struct {
	Name    LibraryName
	Version VersionString
}

func (e MissingVersionError) Error() string {
	return "No version '" + string(e.Version) + "' of library named '" + string(e.Name) + "' found"
}

type VersionConflictError struct {
	Name       LibraryName
	Existing   Version
	Additional Version
}

func (e VersionConflictError) Error() string {
	return "Existing version '" + string(e.Existing.Version) + "' of '" + string(e.Name) +
		"' conflicted with additional version '" + string(e.Additional.Version) + "'"
}

type EmptyLibraryError struct {
	Name LibraryName
}

func (e EmptyLibraryError) Error() string {
	return "No versions associated with library named '" + string(e.Name) + "' found"
}
