package login

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"iceroad/codenight/util"
	"net/http"
)

const userPath = "/user"

var env = util.GetEnv()

var oauthConf = &oauth2.Config{
	ClientID:     env.GithubKey,
	ClientSecret: env.GithubSecret,
	Scopes:       []string{"user:email", "repo"},
	Endpoint:     githuboauth.Endpoint,
}

var oauthStateString = env.GithubStateString

func handleGithubLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGithubCallback(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func AddRoutes(router *httprouter.Router) {
	router.GET("/login", handleGithubLogin)
	router.GET("/oauthCallback", handleGithubCallback)
}
