package login

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"iceroad/codenight/user"
	"iceroad/codenight/util/env"
	"iceroad/codenight/util/session"
	"log"
	"net/http"
)

var environment = env.Get()

var oauthConf = &oauth2.Config{
	ClientID:     environment.GithubKey,
	ClientSecret: environment.GithubSecret,
	Scopes:       []string{"user", "repo"},
	Endpoint:     githuboauth.Endpoint,
}

var oauthStateString = environment.GithubStateString

func loginHandler(res http.ResponseWriter, req *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func oauthCallbackHandler(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if state != oauthStateString {
		log.Printf("ERROR: Invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	code := req.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("ERROR: oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}
	log.Printf("-----TOKEN: %+v \n\n", token)
	encodedToken, _ := json.Marshal(token)

	session.Set(res, req, "github-token", string(encodedToken))

	newUser, userErr := fetchUser(token)
	if userErr != nil {
		log.Printf("ERROR: Failed to fetch user", userErr)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	log.Printf("Upserting User with username '%s'\n", newUser.UserName)

	upsertErr := user.Upsert(newUser)
	if upsertErr != nil {
		log.Printf("ERROR: Failed to create user", upsertErr)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
}

func fetchUser(token *oauth2.Token) (*user.User, error) {
	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	githubUser, _, userError := client.Users.Get("")

	if userError != nil {
		log.Printf("ERROR: client.Users.Get() failed with '%s'\n", userError)
		return nil, userError
	}

	githubEmails, _, emailError := client.Users.ListEmails(&github.ListOptions{Page: 1, PerPage: 10})

	if emailError != nil {
		log.Printf("ERROR: client.Users.ListEmails failed with '%s'\n", emailError)
		return nil, emailError
	}

	user := &user.User{
		Name:      githubUser.Name,
		UserName:  githubUser.Login,
		Email:     selectEmail(githubEmails),
		Blog:      githubUser.Blog,
		Location:  githubUser.Location,
		AvatarUrl: githubUser.AvatarURL}

	return user, nil
}

func selectEmail(emails []*github.UserEmail) *string {
	for i := range emails {
		current := emails[i]
		if *current.Primary {
			return current.Email
		}
	}
	log.Println("ERROR: Failed to find primary email")
	return emails[len(emails)-1].Email
}

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauthCallback", oauthCallbackHandler).Methods(http.MethodGet)
}
