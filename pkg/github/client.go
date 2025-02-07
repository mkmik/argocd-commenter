package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/int128/oauth2-github-app/app"
	"golang.org/x/oauth2"
)

type client struct {
	rest *github.Client
}

func NewClient(ctx context.Context) (*client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return newClientWithPersonalAccessToken(ctx, token), nil
	}
	appID, installationID, privateKey := os.Getenv("GITHUB_APP_ID"), os.Getenv("GITHUB_APP_INSTALLATION_ID"), os.Getenv("GITHUB_APP_PRIVATE_KEY")
	if appID != "" && installationID != "" && privateKey != "" {
		return newClientForGitHubApp(ctx, appID, installationID, privateKey)
	}
	return nil, fmt.Errorf("you need to set either GITHUB_TOKEN or GITHUB_APP_ID")
}

func newClientWithPersonalAccessToken(ctx context.Context, token string) *client {
	c := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	return &client{rest: c}
}

func newClientForGitHubApp(ctx context.Context, appID, installationID, privateKey string) (*client, error) {
	k, err := app.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("invalid GITHUB_APP_PRIVATE_KEY: %w", err)
	}
	cfg := app.Config{
		PrivateKey:     k,
		AppID:          appID,
		InstallationID: installationID,
	}
	c := github.NewClient(oauth2.NewClient(ctx, cfg.TokenSource(ctx)))
	return &client{rest: c}, nil
}
