package user

import (
	"gopkg.in/mgo.v2/bson"
	"iceroad/codenight/db"
)

func GetById(userId string) (User, error) {
	session := db.Connect()
	defer session.Close()
	usersColl := session.DB("codenight").C("users")

	var result User

	err := usersColl.FindId(bson.ObjectIdHex(userId)).One(&result)
	return result, err
}

func CreateIfNotExists(user *User) (string, error) {
	session := db.Connect()
	defer session.Close()
	usersColl := session.DB("codenight").C("users")

	user.Id = bson.NewObjectId()
	err := usersColl.Insert(&user)

	return user.Id.Hex(), err
}
