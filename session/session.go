package session

import (
	"github.com/gorilla/sessions"
	"github.com/rosytucker/codenight/config"
	"log"
	"net/http"
)

// TODO Use redis
var sessionStore *sessions.CookieStore

func SetupStore(environment config.Env) {
	sessionStore = sessions.NewCookieStore([]byte(environment.SessionKey))
}

func Set(res http.ResponseWriter, req *http.Request, key string, value string) {
	session, _ := sessionStore.Get(req, "session")
	session.Values[key] = value
	err := sessions.Save(req, res)
	if err != nil {
		log.Printf("ERROR: could not save sessions: %+v \n", err)
	}
}

func Get(req *http.Request, key string) (string, error) {
	session, err := sessionStore.Get(req, "session")
	if err != nil {
		return "", err
	}
	value := session.Values[key]
	log.Printf("Value from session store: %+v \n", value)
	return value.(string), nil
}
