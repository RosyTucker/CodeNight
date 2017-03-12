package github

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/rosytucker/codenight/config"
	"golang.org/x/oauth2"
	"net/http"
)

type Client struct {
	oauth        oauth2.Config
	githubClient *github.Client
	stateString  string
	Token        oauth2.Token
}

func New(req *http.Request, config ClientConfig) (*Client, error) {
	token, err := getToken(req, config)

	if err != nil {
		return nil, err
	}

	client := &Client{
		oauth:        config.oauth,
		stateString:  config.stateString,
		githubClient: githubClient(config.oauth, token),
		Token:        *token}

	return client, nil
}

func (client *Client) LoginRedirectUrl() string {
	return client.oauth.AuthCodeURL(client.stateString, oauth2.AccessTypeOnline)
}

func (client *Client) GetUser(req *http.Request) (*github.User, error) {

	user, _, userError := client.githubClient.Users.Get("")

	if userError != nil {
		return nil, userError
	}

	emails, _, emailError := client.githubClient.Users.ListEmails(&github.ListOptions{Page: 1, PerPage: 10})

	if emailError != nil {
		config.Log.ErrorF("client.Users.ListEmails failed with '%s'", emailError)
		return nil, emailError
	}

	user.Email = getPrimaryEmail(emails)

	return user, nil
}

func getToken(req *http.Request, config ClientConfig) (*oauth2.Token, error) {
	state := req.FormValue("state")
	if state != config.stateString {
		return nil, errors.Errorf("Invalid oauth state: '%s'\n", state)
	}
	code := req.FormValue("code")
	return config.oauth.Exchange(context.TODO(), code)
}

func githubClient(oauth oauth2.Config, token *oauth2.Token) *github.Client {
	oauthClient := oauth.Client(oauth2.NoContext, token)
	return github.NewClient(oauthClient)
}

func getPrimaryEmail(emails []*github.UserEmail) *string {
	for i := range emails {
		current := emails[i]
		if *current.Primary {
			return current.Email
		}
	}
	config.Log.ErrorF("Failed to find primary email- returning last one added")
	return emails[len(emails)-1].Email
}
