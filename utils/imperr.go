package utils

type MissingLibraryError struct {
	Name LibraryName
}

func (e MissingLibraryError) Error() string {
	return "No library named '"+string(e.Name)+"' found";
}

type MissingVersionError struct {
	Name LibraryName
	Version VersionString
}

func (e MissingVersionError) Error() string {
	return "No version '"+string(e.Version)+"' of library named '"+string(e.Name)+"' found";
}
