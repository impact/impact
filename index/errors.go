package index

import (
	"fmt"
)

type MissingLibraryError struct {
	Name string
}

func (e MissingLibraryError) Error() string {
	return "No library named '" + string(e.Name) + "' found"
}

type MissingVersionError struct {
	Name    string
	Version string
}

func (e MissingVersionError) Error() string {
	return "No version '" + string(e.Version) + "' of library named '" + string(e.Name) + "' found"
}

type VersionConflictError struct {
	Name       string
	Existing   VersionDetails
	Additional VersionDetails
}

func (e VersionConflictError) Error() string {
	return fmt.Sprintf("Existing version '%s'' of '%s' conflicted with additional version '%s'",
		e.Existing.Version.String(), e.Name, e.Additional.Version.String())
}

type EmptyLibraryError struct {
	Name string
}

func (e EmptyLibraryError) Error() string {
	return "No versions associated with library named '" + string(e.Name) + "' found"
}
