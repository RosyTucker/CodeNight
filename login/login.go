package login

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"iceroad/codenight/util"
	"log"
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

func loginHandler(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func oauthCallbackHandler(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	state := req.FormValue("state")
	if state != oauthStateString {
		log.Printf("ERROR: Invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	// code := req.FormValue("code")
	// token, err := oauthConf.Exchange(oauth2.NoContext, code)
	// if err != nil {
	// 	log.Printf("ERROR: oauthConf.Exchange() failed with '%s'\n", err)
	// 	http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
	// 	return
	// }

	sessionCookie := &http.Cookie{Name: "session", Value: "sessionid", HttpOnly: true}
	http.SetCookie(res, sessionCookie)
	http.Redirect(res, req, env.PostLoginRedirect, http.StatusTemporaryRedirect)
}

func AddRoutes(router *httprouter.Router) {
	router.GET("/login", loginHandler)
	router.GET("/oauthCallback", oauthCallbackHandler)
}
