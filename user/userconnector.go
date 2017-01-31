package user

import (
	"github.com/rosytucker/codenight/db"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func GetPublicById(userId string) (*PublicUser, error) {
	user, err := GetById(userId)
	if err != nil {
		return nil, err
	}

	publicUser := &PublicUser{
		Id:          user.Id,
		Name:        user.Name,
		UserName:    user.UserName,
		Description: user.Description,
		Blog:        user.Blog,
		Location:    user.Location,
		AvatarUrl:   user.AvatarUrl}

	return publicUser, err
}

func GetById(userId string) (*User, error) {
	session := db.Connect().Copy()
	defer session.Close()
	usersColl := session.DB("codenight").C("users")

	var result *User
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

func Replace(userId string, updatedUser *User) error {
	session := db.Connect().Copy()
	defer session.Close()
	usersColl := session.DB("codenight").C("users")

	var existingUser User
	err := usersColl.FindId(bson.ObjectIdHex(userId)).One(&existingUser)

	if err != nil {
		return err
	}

	user := &User{
		Id:          existingUser.Id,
		Name:        updatedUser.Name,
		Token:       existingUser.Token,
		UserName:    existingUser.UserName,
		Email:       existingUser.Email,
		Description: updatedUser.Description,
		Blog:        updatedUser.Blog,
		Location:    updatedUser.Location,
		AvatarUrl:   existingUser.AvatarUrl,
		IsAdmin:     existingUser.IsAdmin}

	err = usersColl.UpdateId(existingUser.Id, user)

	return err
}
