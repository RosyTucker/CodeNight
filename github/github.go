package github

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/rosytucker/codenight/config"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"net/http"
)

var oauthConf *oauth2.Config
var oauthStateString string

func ConfigureClient(environment config.Env) {
	oauthConf = &oauth2.Config{
		ClientID:     environment.GithubKey,
		ClientSecret: environment.GithubSecret,
		Scopes:       []string{"user", "repo"},
		Endpoint:     githuboauth.Endpoint}

	oauthStateString = environment.GithubStateString
}

func LoginRedirectUrl() string {
	return oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
}

func GetToken(req *http.Request) (*oauth2.Token, error) {
	state := req.FormValue("state")
	if state != oauthStateString {
		return nil, errors.Errorf("Invalid oauth state: '%s'\n", state)
	}
	code := req.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func GetUser(token *oauth2.Token) (*github.User, error) {
	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	githubUser, _, userError := client.Users.Get("")

	if userError != nil {
		return nil, userError
	}

	githubEmails, _, emailError := client.Users.ListEmails(&github.ListOptions{Page: 1, PerPage: 10})

	if emailError != nil {
		config.Log.ErrorF("client.Users.ListEmails failed with '%s'", emailError)
		return nil, emailError
	}
	githubUser.Email = getPrimaryEmail(githubEmails)

	return githubUser, nil
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
