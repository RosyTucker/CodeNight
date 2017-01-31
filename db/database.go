package db

import (
	"github.com/rosytucker/codenight/config"
	"gopkg.in/mgo.v2"
	"log"
)

var session *mgo.Session

func EstablishInitialConnection(environment config.Env) {
	log.Println("---- Connecting to DB ----")

	var err error
	session, err = mgo.Dial(environment.MongoConnectionString)

	if err != nil {
		log.Printf("ERROR: Can't connect to mongo, go error %+v\n", err)
		panic(err)
	}

	session.SetSafe(&mgo.Safe{})
}

func Connect() *mgo.Session {
	return session.Copy()
}
