package dirinfo

type LocalLibrary struct {
	Name      string // Name of library
	Path      string // Path to library (relative to directory/impact.json file)
	IsFile    bool   // If the library is stored as a single file
	IssuesURL string
}

type DirectoryInfo struct {
	Author    string
	Libraries []LocalLibrary
}
