package session

import (
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
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
		log.Printf("ERROR: Could not get session from request %+v \n", err)
		return "", err
	}
	value := session.Values[key]

	if value == nil {
		return "", errors.Errorf("Value for key %s not found in session store", key)
	}
	return value.(string), nil
}
