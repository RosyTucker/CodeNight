package user

import (
	"gopkg.in/mgo.v2/bson"
	"iceroad/codenight/db"
	"log"
)

func GetById(userId string) (User, error) {
	session := db.Connect().Copy()
	defer session.Close()
	usersColl := session.DB("codenight").C("users")

	var result User

	err := usersColl.FindId(bson.ObjectIdHex(userId)).One(&result)
	return result, err
}

func CreateIfNotExists(user *User) (string, error) {
	session := db.Connect().Copy()
	defer session.Close()
	usersColl := session.DB("codenight").C("users")

	var existingUser *User
	err := usersColl.Find(bson.M{"username": user.UserName}).One(&existingUser)

	if err != nil && err.Error() != "not found" {
		return "", err
	}

	if existingUser != nil {
		log.Printf("User with username '%s' already exists, not creating another\n", user.UserName)
		return existingUser.Id.Hex(), nil
	}

	log.Printf("User with username '%s' does not exist, creating new user\n", user.UserName)
	user.Id = bson.NewObjectId()
	err = usersColl.Insert(&user)

	return user.Id.Hex(), err
}
