package gh

import (
	"context"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func NewClient(githubToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
