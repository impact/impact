package index

import (
	"fmt"
)

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
	Existing   VersionDetails
	Additional VersionDetails
}

func (e VersionConflictError) Error() string {
	return fmt.Sprintf("Existing version '%s'' of '%s' conflicted with additional version '%s'",
		e.Existing.Version.String(), e.Name, e.Additional.Version.String())
}

type EmptyLibraryError struct {
	Name LibraryName
}

func (e EmptyLibraryError) Error() string {
	return "No versions associated with library named '" + string(e.Name) + "' found"
}
