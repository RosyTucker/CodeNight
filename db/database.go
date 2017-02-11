package db

import (
	"github.com/rosytucker/codenight/config"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func EstablishInitialConnection(environment config.Env) {
	config.Log.Debug("---- Connecting to DB ----")

	var err error
	session, err = mgo.Dial(environment.MongoConnectionString)

	if err != nil {
		config.Log.PanicF("Can't connect to mongo '%+v'", err)
		panic(err)
	}

	session.SetSafe(&mgo.Safe{})
}

func Connect() *mgo.Session {
	return session.Copy()
}
