package db

import (
	"gopkg.in/mgo.v2"
	"iceroad/codenight/env"
	"log"
)

var session *mgo.Session

func init() {
	log.Println("---- Connecting to DB ----")
	environment := env.Get()

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
