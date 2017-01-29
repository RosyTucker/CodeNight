package session

import (
	"github.com/gorilla/sessions"
	"iceroad/codenight/env"
	"log"
	"net/http"
)

// TODO Use redis
var sessionStore *sessions.CookieStore

func init() {
	var env = env.Get()
	sessionStore = sessions.NewCookieStore([]byte(env.SessionKey))
}

func Set(res http.ResponseWriter, req *http.Request, key string, value string) {
	session, _ := sessionStore.Get(req, "session")
	session.Values[key] = value
	err := sessions.Save(req, res)
	if err != nil {
		log.Printf("ERROR: could not save sessions: %+v \n", err)
	}
}

func Get(req *http.Request, key string) string {
	session, _ := sessionStore.Get(req, "session")
	value := session.Values[key]
	return value.(string)
}
