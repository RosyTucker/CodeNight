package github

import (
	"github.com/rosytucker/codenight/config"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type ClientConfig struct {
	oauth            oauth2.Config
	stateString      string
	LoginRedirectUrl string
}

func CreateConfig(environment config.Env) ClientConfig {
	oauthConfig := oauth2.Config{
		ClientID:     environment.GithubKey,
		ClientSecret: environment.GithubSecret,
		Scopes:       []string{"user"},
		Endpoint:     githuboauth.Endpoint}

	stateString := environment.GithubStateString

	return ClientConfig{
		oauth:            oauthConfig,
		stateString:      stateString,
		LoginRedirectUrl: oauthConfig.AuthCodeURL(stateString, oauth2.AccessTypeOnline)}
}