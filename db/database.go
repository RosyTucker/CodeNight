package db

import (
	"gopkg.in/mgo.v2"
	"iceroad/codenight/env"
	"log"
)

func Connect() (session *mgo.Session) {
	environment := env.Get()
	session, err := mgo.Dial(environment.MongoConnectionString)
	if err != nil {
		log.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	session.SetSafe(&mgo.Safe{})
	return session
}
