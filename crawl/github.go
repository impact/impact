package crawl

import (
	"fmt"
	"log"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"

	"github.com/xogeny/impact/parsing"
	"github.com/xogeny/impact/recorder"
)

type GitHubCrawler struct {
	token string
	user  string
}

var exclusionList []string

func init() {
	exclusionList = []string{
		"modelica-3rdparty:BrineProp:0.1.9",      // Directory structure is a mess
		"modelica-3rdparty:BondLib:2.3",          // Too large
		"modelica-3rdparty:DESLib:1.6.1",         // Too large
		"modelica-3rdparty:FCSys:0.2.3",          // Tag dir mismatch
		"modelica-3rdparty:FCSys:0.2.2",          // Tag dir mismatch
		"modelica-3rdparty:FCSys:0.2.1",          // Tag dir mismatch
		"modelica-3rdparty:FCSys:0.2.0",          // Tag dir mismatch
		"modelica-3rdparty:FCSys:0.1.2",          // Tag dir mismatch
		"modelica-3rdparty:FCSys:0.1.1",          // Tag dir mismatch
		"modelica-3rdparty:FCSys:0.1.0",          // Tag dir mismatch
		"modelica-3rdparty:HelmholtzMedia:0.9.1", // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.9.0", // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.8.4", // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.8.2", // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.8.1", // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.8",   // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.7.1", // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.7",   // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.6",   // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.6.1", // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.5",   // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.4",   // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.3",   // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.2",   // Unsupported directory structure
		"modelica-3rdparty:HelmholtzMedia:0.1",   // Unsupported directory structure
		"modelica-3rdparty:LinearMPC:0.1",        // Tag dir mismatch
		"modelica-3rdparty:ModelicaDEVS:1.0",     // Self reference (and invalid at that)
		"modelica-3rdparty:MotorcycleLib:1.0",    // Too large
		"modelica-3rdparty:NCLib:0.82",           // Missing package.mo
		"modelica:Modelica_LinearSystems2:2.3.1", // Dir name error
	}
}

func exclude(user string, reponame string, tagname string) bool {
	str := fmt.Sprintf("%s:%s:%s", user, reponame, tagname)
	for _, ex := range exclusionList {
		//log.Printf("Comparing '%s' to '%s'", ex, str)
		if ex == str {
			return true
		}
	}
	return false
}

func (c GitHubCrawler) Crawl(r recorder.Recorder, verbose bool, logger *log.Logger) error {
	// Start with whatever token we were given when this crawler was created
	token := c.token

	// If a token wasn't provided with the crawler, look for a token
	// as an environment variable
	if c.token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	// Create a client assuming no authentication
	client := github.NewClient(nil)

	// If we have a token, re-initialize the client with
	// authentication
	if token != "" {
		tok := &oauth.Transport{
			Token: &oauth.Token{AccessToken: token},
		}
		client = github.NewClient(tok.Client())
	}

	// Get a list of all repositories associated with the specified
	// organization
	repos, _, err := client.Repositories.List(c.user, nil)
	if err != nil {
		logger.Printf("Error listing repositories for %s: %v", c.user, err)
		return fmt.Errorf("Error listing repositories for %s: %v", c.user, err)
	}

	// Loop over all repos associated with the given owner
	for _, repo := range repos {
		if verbose {
			logger.Printf("Processing: %s (%s)", *repo.Name, *repo.HTMLURL)
		}

		// Get all the tags associated with this repository
		tags, _, err := client.Repositories.ListTags(c.user, *repo.Name, nil)
		if err != nil {
			logger.Printf("Error getting tags for repository %s/%s: %v",
				c.user, *repo.Name, err)
			continue
		}

		// Loop over the tags
		for _, tag := range tags {
			// Check if this has a semantic version
			name := *tag.Name
			if name[0] == 'v' {
				name = name[1:]
			}

			// Check for version we know are not supported
			if exclude(c.user, *repo.Name, name) {
				continue
			}

			v, verr := parsing.NormalizeVersion(name)
			if verr != nil {
				// If not, ignore it
				if verbose {
					logger.Printf("  %s: Ignoring", name)
				}
				continue
			}

			if verbose {
				logger.Printf("  %s: Recording", name)
			}

			// Formulate directory info (impact.json) for this version of this repository
			di := ExtractInfo(client, c.user, repo, tag, name, logger)

			if len(di.Libraries) == 0 {
				logger.Printf("    No Modelica libraries found in repository %s:%s",
					*repo.Name, name)
				continue
			}

			// Loop over all libraries present in this repository
			for _, lib := range di.Libraries {
				if verbose {
					logger.Printf("    Processing library %s @ %s", lib.Name, lib.Path)
				}
				libr := r.GetLibrary(di.Owner, lib.Name)

				if repo.Description != nil {
					libr.SetDescription(*repo.Description)
				}
				if repo.HTMLURL != nil {
					libr.SetHomepage(*repo.HTMLURL)
				}
				libr.SetStars(*repo.StargazersCount)
				if repo.Owner.Email != nil {
					libr.SetEmail(*repo.Owner.Email)
				}

				vr := libr.AddVersion(v)
				vr.SetHash(*tag.Commit.SHA)

				for _, dep := range lib.Dependencies {
					vr.AddDependency(dep.Name, dep.Version)
				}
			}
		}
	}
	return nil
}

func MakeGitHubCrawler(user string, token string) GitHubCrawler {
	return GitHubCrawler{
		token: token,
		user:  user,
	}
}

var _ Crawler = (*GitHubCrawler)(nil)
