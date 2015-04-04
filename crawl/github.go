package crawl

import (
	"log"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

type GitHubCrawler struct {
	token string
	user  string
}

func (c GitHubCrawler) Crawl() error {
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
		return err
	}

	for _, repo := range repos {
		log.Printf("repo = %v", repo)
	}
	return nil
}

func MakeGitHubCrawler(user string, token string) GitHubCrawler {
	return GitHubCrawler{
		token: token,
		user:  user,
	}
}
