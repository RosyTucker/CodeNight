package login

import (
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"iceroad/codenight/session"
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
	http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
}

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauthCallback", oauthCallbackHandler).Methods(http.MethodGet)
}
