package util

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

func init() {
	var env = GetEnv()
	session, err := mgo.Dial(env.MongoConnectionString)

	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
}

func Users() *mgo.Collection {
	return session.DB("codenight").C("users")
}
