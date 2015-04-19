package crawl

import (
	"fmt"
	"log"

	"github.com/blang/semver"

	"github.com/google/go-github/github"

	"github.com/xogeny/impact/dirinfo"
	"github.com/xogeny/impact/parsing"
)

func getUses(client *github.Client, user string, reponame string,
	path string, opts *github.RepositoryContentGetOptions) (map[string]semver.Version, error) {
	blank := map[string]semver.Version{}

	// Read contents of top-level package
	lcon, _, _, err := client.Repositories.GetContents(user, reponame, path, opts)
	if err != nil {
		return blank,
			fmt.Errorf("Error while reading contents of %s in github repository %s: %v",
				path, reponame, err)
	}

	dec, err := lcon.Decode()
	if err != nil {
		return blank,
			fmt.Errorf("Error while decoding contents of %s in github repository %s: %v",
				path, reponame, err)
	}

	contents := string(dec)

	uses, err := parsing.ParseUses(contents)
	if err != nil {
		return blank,
			fmt.Errorf("Error while parsing contents of %s in github repository %s: %v",
				path, reponame, err)
	}

	return uses, nil
}

func Exists(client *github.Client, user string, reponame string,
	path string, opts *github.RepositoryContentGetOptions) (file bool, dir bool) {
	f, d, _, err := client.Repositories.GetContents(user, reponame, path, opts)
	if err != nil {
		return false, false
	}
	if f != nil {
		return true, false
	}
	if d != nil {
		return false, true
	}
	return false, false
}

// The goal of this function is to construct a DirectoryInfo object.  It does this by first
// reading whatever directory information it can find in impact.json.  Then it tries to
// "infer" the rest using some heuristics (to lower the burden on library developers)
func ExtractInfo(client *github.Client, user string, altname string, repo github.Repository,
	sha string, versionString string, verbose bool, logger *log.Logger) dirinfo.DirectoryInfo {

	// Extract the name of the respository
	repostr := *repo.Name

	// Extract information about the owner of this repository (note: this is a URI)
	owner_uri := user
	if repo.Owner.HTMLURL != nil {
		owner_uri = *repo.Owner.HTMLURL
	}

	// Get the owner's email address, if provided
	email := ""
	if repo.Owner.Email != nil {
		email = *repo.Owner.Email
	}

	// Specify which version of the repository we are interested in
	opts := &github.RepositoryContentGetOptions{
		Ref: sha,
	}

	// Create a "blank" directory info as default
	di := dirinfo.MakeDirectoryInfo()

	// Parse any impact.json file the is present
	dcon, _, _, err := client.Repositories.GetContents(user, repostr, "impact.json", opts)

	// If impact.json exists, parse it and use that as our baseline
	if dcon != nil && err == nil {
		pdi, perr := dirinfo.Parse(*dcon.Content)
		if perr == nil {
			log.Printf("Parsed impact.json file in %s: %v", repostr, pdi)
			di = pdi
		} else {
			log.Printf("Unable to parse impact.json in %s: %v", repostr, perr)
		}
	}

	// NOW, use heuristics to infer missing information

	// Is owner information missing.  If so, use either the owner whose repositories
	// are being scanned or the owner of this particular repository (see logic above
	// regarding 'owner')
	if di.OwnerURI == "" {
		di.OwnerURI = owner_uri
	}

	// Ditto with email
	if di.Email == "" {
		di.Email = email
	}

	// This is local function where only the path is an argument, all the other
	// arguments come from the local environment.
	lookupPath := func(path string) (bool, bool) {
		return Exists(client, user, repostr, path, opts)
	}

	// Are any libraries mentioned?  If not, we need to figure out what the structure
	// is here.  There are two patterns.  Either a file named <RepoName>.mo or a
	// directory named <RepoName>.  If neither of these conventions is followed, the
	// library developers needs to add an explicit impact.json
	if len(di.Libraries) == 0 {
		// Provide default values for path and isfile
		path := ""
		isfile := true

		// This is a simplified version of the version string that strips build
		// information away.  This is necessary because some library developers
		// don't include build information in the names of files and directories
		// that include an explicit version number.  So we check for this as well.
		sver := parsing.SimpleVersion(versionString)

		// Flag to indicate whether we have found a library
		found := false

		// A list of Modelica files to look for that could be libraries
		// (stored as files)
		filenames := []string{
			repostr + ".mo",
			repostr + " " + versionString + ".mo",
			altname + ".mo",
			altname + " " + versionString + ".mo",
		}

		// A list of directories to look for that could be Modelica libraries
		// (stored as directories)
		dirnames := []string{
			repostr,
			repostr + " " + versionString,
			altname,
			altname + " " + versionString,
		}

		// If the simple version string isn't the same as the full version string,
		// add a few more variations.
		if sver != versionString {
			filenames = append(filenames, repostr+" "+sver+".mo")
			filenames = append(filenames, altname+" "+sver+".mo")
			dirnames = append(dirnames, repostr+" "+sver)
			dirnames = append(dirnames, altname+" "+sver)
		}

		// First, check for the library stored as a file...
		for _, fpath := range filenames {
			// Look for that file in the repository
			fexists, _ := lookupPath(fpath)
			if verbose {
				log.Printf("Looking for file %s in %s:%s: %v", fpath, repostr,
					versionString, fexists)
			}
			// If no library has been found yet and this file exists, record
			// its information and set the found flag
			if !found && fexists {
				path = fpath
				found = true
				if verbose {
					log.Printf("  Found")
				}
			}
		}

		// Next, check for the library stored as a directory...
		for _, dpath := range dirnames {
			// Look for a directory in the repository with the matching name
			// TODO: Look for package.mo inside!
			_, dexists := lookupPath(dpath)
			if verbose {
				log.Printf("Looking for directory %s in %s:%s: %v", dpath, repostr,
					versionString, dexists)
			}
			// If we haven't yet found a match and this directory exists, then
			// record the information and set found flag
			if !found && dexists {
				path = dpath
				isfile = false
				found = true
				if verbose {
					log.Printf("  Found")
				}
			}
		}

		// Finally, if we still haven't found a library, chec to see if the
		// repository itself is a library
		if !found {
			// Check if repository IS a package
			fexists, _ := lookupPath("package.mo")
			if verbose {
				log.Printf("Check if directory %s:%s is a package: %v", repostr, versionString,
					fexists)
			}
			// If so, record information about it
			if fexists {
				path = ""
				isfile = false
				found = true
				if verbose {
					log.Printf("  Found")
				}
			}
		}

		// If we found a library, create an instance of dirinfo.LocalLibrary
		// N.B. - These heuristics only look for one library in the repository.
		// If a library developer wants to store more than one, they need to
		// state this explicitly in the impact.json file and not rely on us
		// inferring it somehow.
		if found {
			di.Libraries = []*dirinfo.LocalLibrary{
				&dirinfo.LocalLibrary{
					Name:         repostr,
					Path:         path,
					IsFile:       isfile,
					Dependencies: []dirinfo.Dependency{},
				},
			}
		} else {
			log.Printf("Nothing found in %s/%s for %v or %v", repostr, filenames, dirnames)
		}
	}

	// Now, let's loop over all the libraries we are aware of...
	for _, lib := range di.Libraries {
		// Determine path to top-level package in repository
		path := lib.Path
		if !lib.IsFile {
			path = fmt.Sprintf("%s/package.mo", lib.Path)
		}

		// Extract information about any libraries this library uses
		uses, err := getUses(client, user, repostr, path, opts)
		if err != nil {
			log.Printf("Error extracting uses annotation: %v", err)
			continue
		}

		for libname, ver := range uses {
			lib.Dependencies = append(lib.Dependencies, dirinfo.Dependency{
				Name:    libname,
				Version: ver,
			})
		}

		if lib.IssuesURL == "" {
			lib.IssuesURL = *repo.IssuesURL
		}
	}

	return di
}
