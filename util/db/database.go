package db

import (
	"gopkg.in/mgo.v2"
	"iceroad/codenight/util/env"
	"log"
)

// var session *mgo.Session
// var db *mgo.Database

// func init() {
// 	var environment = env.Get()
// 	log.Printf(environment.MongoConnectionString)
// 	session, err := mgo.Dial(environment.MongoConnectionString)

// 	if err != nil {
// 		log.Fatal("ERROR: Could connect to DB: %+v\n", err)
// 		panic(err)
// 	} else {
// 		log.Printf("SUCCESS: Connected to DB")
// 	}

// 	defer session.Close()

// 	session.SetMode(mgo.Monotonic, true)
// 	session.SetSafe(&mgo.Safe{})
// 	db := session.DB("codenight")
// 	configureUsers(db)
// }

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

// func configureUsers(database *mgo.Database) {
// 	index := mgo.Index{
// 		Key:    []string{"id"},
// 		Unique: true,
// 	}
// 	err := database.C("users").EnsureIndex(index)

// 	if err != nil {
// 		log.Fatal("Failed to configure users collection")
// 		panic(err)
// 	}
// }
