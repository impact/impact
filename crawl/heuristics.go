package crawl

import (
	"github.com/google/go-github/github"
	"log"

	"github.com/xogeny/impact/dirinfo"
)

func ExtractInfo(client *github.Client, user string, repo github.Repository,
	tag github.RepositoryTag, logger *log.Logger) dirinfo.DirectoryInfo {

	// TODO: What if this is a mirror?  Follow "Source" repository?

	repostr := *repo.Name

	owner := user
	if repo.Owner.Login != nil {
		owner = *repo.Owner.Login
	}
	ref := *tag.Commit.SHA
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// Parse any impact.json file the is present
	dcon, _, _, err := client.Repositories.GetContents(owner, repostr, "impact.json", opts)

	// Create a "blank" directory info as default
	di := dirinfo.MakeDirectoryInfo()

	// If impact.json exists, parse it and use that as our baseline
	if dcon != nil && err == nil {
		di = dirinfo.Parse(*dcon.Content)
	}

	// NOW, use heuristics to infer missing information

	// Is owner information missing.  If so, use either the owner whose repositories
	// are being scanned or the owner of this particular repository (see logic above
	// regarding 'owner')
	if di.Owner == "" {
		di.Owner = owner
	}

	// Are any libraries mentioned?  If not, we need to figure out what the structure
	// is here.  There are two patterns.  Either a file named <RepoName>.mo or a
	// directory named <RepoName>.  If neither of these conventions is followed, the
	// library developers needs to add an explicit impact.json
	if len(di.Libraries) == 0 {
		lcon, _, _, err := client.Repositories.GetContents(owner, repostr, repostr+".mo", opts)
		if lcon != nil && err == nil {
			di.Libraries = []dirinfo.LocalLibrary{
				dirinfo.LocalLibrary{
					Name:   repostr,
					Path:   repostr + ".mo",
					IsFile: true,
				},
			}
		} else {
			_, ldcon, _, err := client.Repositories.GetContents(owner, repostr, repostr, opts)
			if ldcon != nil && err == nil {
				// If directory named RepoName exists...
				di.Libraries = []dirinfo.LocalLibrary{
					dirinfo.LocalLibrary{
						Name:   repostr,
						Path:   repostr,
						IsFile: false,
					},
				}
			}
		}
	}

	for _, lib := range di.Libraries {
		if lib.IssuesURL == "" {
			lib.IssuesURL = *repo.IssuesURL
		}
	}

	return di
}
