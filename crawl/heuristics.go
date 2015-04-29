package crawl

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/blang/semver"

	"github.com/google/go-github/github"

	"github.com/xogeny/impact/dirinfo"
	"github.com/xogeny/impact/parsing"
)

func parsePackage(client *github.Client, user string, reponame string,
	mopath string, opts *github.RepositoryContentGetOptions) (string,
	map[string]semver.Version, error) {
	blank := map[string]semver.Version{}

	reader, err := client.Repositories.DownloadContents(user, reponame, mopath, opts)
	if err != nil {
		return "", blank, fmt.Errorf("Unable to download Modelica code for %s: %v", mopath)
	}
	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", blank, fmt.Errorf("Error reading response: %v", err)
	}

	contents := string(raw)

	uses, err := parsing.ParseUses(contents)
	if err != nil {
		return "", blank,
			fmt.Errorf("Error while parsing uses annotation of %s in github repository %s: %v",
				mopath, reponame, err)
	}

	name, err := parsing.ParseName(contents)
	if err != nil {
		return "", blank,
			fmt.Errorf("Error while parsing name of %s in github repository %s: %v",
				mopath, reponame, err)
	}

	return name, uses, nil
}

func getName(code io.ReadCloser) (string, error) {
	return "", fmt.Errorf("Unimplemented")
}

func getLibraries(client *github.Client, user string, repostr string, verbose bool,
	opts *github.RepositoryContentGetOptions) ([]*dirinfo.LocalLibrary, error) {
	blank := []*dirinfo.LocalLibrary{}

	if verbose {
		log.Printf("  Reviewing contents of %s/%s", user, repostr)
	}
	// Grab information about the contents of repository's root directory
	_, dcon, _, err := client.Repositories.GetContents(user, repostr, ".", opts)
	if err != nil {
		return blank, fmt.Errorf("Unable to fetch repository files: %v", err)
	}

	// First check to see if the root of the repository contains a package.mo
	// file.  If so, the whole repository is a library...
	for _, con := range dcon {
		if *con.Name == "package.mo" {
			if verbose {
				log.Printf("  Repository is a library")
			}
			body, err := client.Repositories.DownloadContents(user, repostr, *con.Path, opts)
			if err != nil {
				log.Printf("Unable to read contents of %s/%s/%s: %v", user, repostr, *con.Path,
					err)
				continue
			}
			name, err := getName(body)
			if err != nil {
				log.Printf("Unable to extract name from %s/%s/%s", user, repostr, *con.Path)
				continue
			}
			return []*dirinfo.LocalLibrary{
				&dirinfo.LocalLibrary{
					Name:         name,
					Path:         ".",
					IsFile:       false,
					Dependencies: []dirinfo.Dependency{},
				},
			}, nil
		}
	}

	ret := []*dirinfo.LocalLibrary{}
	for _, con := range dcon {
		switch *con.Type {
		case "file":
			if strings.HasSuffix(*con.Name, ".mo") {
				ret = append(ret, &dirinfo.LocalLibrary{
					Name:         repostr,
					Path:         *con.Path,
					IsFile:       true,
					Dependencies: []dirinfo.Dependency{},
				})
			}
		case "dir":
			_, subcons, _, err := client.Repositories.GetContents(user, repostr, *con.Name, opts)
			if err != nil {
				continue
			}
			for _, sub := range subcons {
				if *sub.Name == "package.mo" {
					ret = append(ret, &dirinfo.LocalLibrary{
						Name:         repostr,
						Path:         *con.Path,
						IsFile:       false,
						Dependencies: []dirinfo.Dependency{},
					})
				}
			}
		}
	}

	return blank, fmt.Errorf("No libraries found in %s/%s", user, repostr)
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
	fcon, _, _, err := client.Repositories.GetContents(user, repostr, "impact.json", opts)

	// If impact.json exists, parse it and use that as our baseline
	if fcon != nil && err == nil {
		pdi, perr := dirinfo.Parse(*fcon.Content)
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

	// Are any libraries mentioned?  If not, we need to figure out what the structure
	// is here.  There are two patterns.  Either a file named <RepoName>.mo or a
	// directory named <RepoName>.  If neither of these conventions is followed, the
	// library developers needs to add an explicit impact.json
	if len(di.Libraries) == 0 {
		libs, err := getLibraries(client, user, repostr, verbose, opts)
		if err != nil {
			if verbose {
				log.Printf("No libraries found in %s/%s", user, repostr)
			}
		}
		di.Libraries = libs
	}

	// Now, let's loop over all the libraries we are aware of...
	for _, lib := range di.Libraries {
		// Determine path to top-level package in repository
		path := lib.Path
		if !lib.IsFile {
			path = fmt.Sprintf("%s/package.mo", lib.Path)
		}

		// Extract information about any libraries this library uses
		name, uses, err := parsePackage(client, user, repostr, path, opts)
		if err != nil {
			log.Printf("Error extracting uses annotation: %v", err)
			continue
		}

		lib.Name = name

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
