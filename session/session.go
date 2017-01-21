package session

import (
	"github.com/gorilla/sessions"
	"iceroad/codenight/util"
	"log"
	"net/http"
)

var env = util.GetEnv()

// TODO Use redis
var sessionStore = sessions.NewCookieStore([]byte(env.SessionKey))

func Set(res http.ResponseWriter, req *http.Request, key string, value string) {
	session, _ := sessionStore.Get(req, "session")
	session.Values[key] = value
	err := sessions.Save(req, res)
	if err != nil {
		log.Printf("ERROR: could not save sessions: %+v \n", err)
	}
}

func Get(req *http.Request, key string) interface{} {
	session, _ := sessionStore.Get(req, "session")
	value := session.Values[key]
	return value
}
