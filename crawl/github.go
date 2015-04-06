package crawl

import (
	"fmt"
	"log"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/blang/semver"
	"github.com/google/go-github/github"

	"github.com/xogeny/impact/recorder"
)

type GitHubCrawler struct {
	token string
	user  string
}

func (c GitHubCrawler) Crawl(r recorder.Recorder, logger *log.Logger) error {
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
		logger.Printf("Processing: %s (%s)", *repo.Name, *repo.HTMLURL)

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
			v, verr := semver.Parse(name)
			if verr != nil {
				// If not, ignore it
				logger.Printf("  %s: Ignoring", name)
				continue
			}

			logger.Printf("  %s: Recording", name)

			// Formulate directory info (impact.json) for this version of this repository
			di := ExtractInfo(client, c.user, repo, tag, logger)
			if len(di.Libraries) == 0 {
				logger.Printf("    No Modelica libraries found in repository %s", *repo.Name)
				continue
			}

			// Loop over all libraries present in this repository
			for _, lib := range di.Libraries {
				logger.Printf("    Processing library %s @ %s", lib.Name, lib.Path)
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
