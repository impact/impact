package crawl

import (
	"fmt"
	"log"

	"github.com/blang/semver"

	"github.com/google/go-github/github"

	"github.com/xogeny/impact/dirinfo"
	"github.com/xogeny/impact/parsing"
)

func getUses(client *github.Client, owner string, reponame string,
	path string, opts *github.RepositoryContentGetOptions) (map[string]semver.Version, error) {
	blank := map[string]semver.Version{}

	// Read contents of top-level package
	lcon, _, _, err := client.Repositories.GetContents(owner, reponame, path, opts)
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

func Exists(client *github.Client, owner string, reponame string,
	path string, opts *github.RepositoryContentGetOptions) (file bool, dir bool) {
	f, d, _, err := client.Repositories.GetContents(owner, reponame, path, opts)
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

func ExtractInfo(client *github.Client, user string, repo github.Repository,
	sha string, versionString string, logger *log.Logger) dirinfo.DirectoryInfo {

	// TODO: What if this is a mirror?  Follow "Source" repository?

	repostr := *repo.Name

	owner := user
	if repo.Owner.Login != nil {
		owner = *repo.Owner.Login
	}

	opts := &github.RepositoryContentGetOptions{
		Ref: sha,
	}

	// Create a "blank" directory info as default
	di := dirinfo.MakeDirectoryInfo()

	// Parse any impact.json file the is present
	dcon, _, _, err := client.Repositories.GetContents(owner, repostr, "impact.json", opts)

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
	if di.Owner == "" {
		di.Owner = owner
	}

	lookupPath := func(path string) (bool, bool) {
		return Exists(client, owner, repostr, path, opts)
	}

	// Are any libraries mentioned?  If not, we need to figure out what the structure
	// is here.  There are two patterns.  Either a file named <RepoName>.mo or a
	// directory named <RepoName>.  If neither of these conventions is followed, the
	// library developers needs to add an explicit impact.json
	if len(di.Libraries) == 0 {
		path := ""
		isfile := true
		found := false
		sver := parsing.SimpleVersion(versionString)

		filenames := []string{repostr + ".mo", repostr + " " + versionString + ".mo"}
		dirnames := []string{repostr, repostr + " " + versionString}

		if sver != versionString {
			filenames = append(filenames, repostr+" "+sver+".mo")
			dirnames = append(dirnames, repostr+" "+sver)
		}

		for _, fpath := range filenames {
			//log.Printf("Looking for file %s in %s:%s", fpath, repostr, versionString)
			fexists, _ := lookupPath(fpath)
			if !found && fexists {
				path = fpath
				found = true
			}
		}

		for _, dpath := range dirnames {
			//log.Printf("Looking for directory %s in %s:%s", dpath, repostr, versionString)
			_, dexists := lookupPath(dpath)
			if !found && dexists {
				path = dpath
				isfile = false
				found = true
			}
		}

		if !found {
			// Check if repository IS a package
			if fexists, _ := lookupPath("package.mo"); fexists {
				path = ""
				isfile = false
				found = true
			}
		}

		if found {
			di.Libraries = []*dirinfo.LocalLibrary{
				&dirinfo.LocalLibrary{
					Name:         repostr,
					Path:         path,
					IsFile:       isfile,
					Dependencies: []dirinfo.Dependency{},
				},
			}
		}
	}

	for _, lib := range di.Libraries {
		// Determine path to top-level package in repository
		path := lib.Path
		if !lib.IsFile {
			path = fmt.Sprintf("%s/package.mo", lib.Path)
		}

		// Read contents of top-level package
		lcon, _, _, err := client.Repositories.GetContents(owner, repostr, path, opts)
		if err != nil {
			log.Printf("Error while reading contents of %s in github repository %s: %v",
				path, repostr, err)
			continue
		}

		dec, err := lcon.Decode()
		if err != nil {
			log.Printf("Error while decoding contents of %s in github repository %s: %v",
				path, repostr, err)
			continue
		}

		contents := string(dec)

		uses, err := parsing.ParseUses(contents)
		if err != nil {
			log.Printf("Error while parsing contents of %s in github repository %s: %v",
				path, repostr, err)
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
