package login

import (
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"iceroad/codenight/session"
	"iceroad/codenight/user"
	"iceroad/codenight/util"
	"log"
	"net/http"
)

var env = util.GetEnv()

var oauthConf = &oauth2.Config{
	ClientID:     env.GithubKey,
	ClientSecret: env.GithubSecret,
	Scopes:       []string{"user:email", "repo"},
	Endpoint:     githuboauth.Endpoint,
}

var oauthStateString = env.GithubStateString

func loginHandler(res http.ResponseWriter, req *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func oauthCallbackHandler(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if state != oauthStateString {
		log.Printf("ERROR: Invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	code := req.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("ERROR: oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}
	session.Set(res, req, "github-token", token.AccessToken)
	user, userErr := fetchUser(token)
	if userErr != nil {
		log.Printf("ERROR: Failed to fetch user", userErr)
		http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	log.Printf("Upserting User with username '%s'\n", *user.UserName)

	upsertErr := user.Upsert(githubUser, token)
	if upsertErr != nil {
		log.Printf("ERROR: Failed to create user", upsertErr)
		http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
}

func fetchUser(token *oauth2.Token) (*user.User, error) {
	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	githubUser, _, err := client.Users.Get("")
	if err != nil {
		log.Printf("ERROR: client.Users.Get() failed with '%s'\n", err)
		return nil, err
	}
	user := &user.User{
		Id:        githubUser.ID,
		Name:      githubUser.Name,
		UserName:  githubUser.Login,
		Email:     githubUser.Email,
		Blog:      githubUser.Blog,
		Location:  githubUser.Location,
		AvatarUrl: githubUser.AvatarURL}

	return user, nil
}

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauthCallback", oauthCallbackHandler).Methods(http.MethodGet)
}
