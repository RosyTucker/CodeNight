package user

import (
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/db"
	"gopkg.in/mgo.v2/bson"
)

var dbName = config.GetEnv().DbName

func GetPublicById(userId string) (*PublicUser, error) {
	user, err := GetById(userId)
	if err != nil {
		return nil, err
	}

	publicUser := &PublicUser{
		Id:          user.Id,
		Name:        user.Name,
		Description: user.Description,
		Blog:        user.Blog,
		Location:    user.Location,
		AvatarUrl:   user.AvatarUrl}

	return publicUser, err
}

func GetById(userId string) (*User, error) {
	session := db.Connect().Copy()
	defer session.Close()
	usersColl := session.DB(dbName).C("users")

	var result *User
	err := usersColl.FindId(bson.ObjectIdHex(userId)).One(&result)
	return result, err
}

func CreateIfNotExists(user *User) (string, error) {
	session := db.Connect().Copy()
	defer session.Close()
	usersColl := session.DB(dbName).C("users")

	var existingUser *User
	err := usersColl.Find(bson.M{"username": user.UserName}).One(&existingUser)

	if err != nil && err.Error() != "not found" {
		return "", err
	}

	if existingUser != nil {
		config.Log.InfoF("User with username '%s' already exists, not creating another", user.UserName)
		return existingUser.Id.Hex(), nil
	}

	config.Log.InfoF("User with username '%s' does not exist, creating new user", user.UserName)
	user.Id = bson.NewObjectId()
	err = usersColl.Insert(&user)

	return user.Id.Hex(), err
}

func Replace(userId string, updatedUser *PublicUser) error {
	session := db.Connect().Copy()
	defer session.Close()
	dbCollection := session.DB(dbName).C("users")

	var existingUser User
	err := dbCollection.FindId(bson.ObjectIdHex(userId)).One(&existingUser)

	if err != nil {
		return err
	}

	update := &User{
		Id:          existingUser.Id,
		Name:        updatedUser.Name,
		Token:       existingUser.Token,
		UserName:    existingUser.UserName,
		Email:       existingUser.Email,
		Description: updatedUser.Description,
		Blog:        updatedUser.Blog,
		Location:    updatedUser.Location,
		AvatarUrl:   existingUser.AvatarUrl,
		MemberSince: existingUser.MemberSince,
		Company:     updatedUser.Company,
		IsAdmin:     existingUser.IsAdmin}

	err = dbCollection.UpdateId(existingUser.Id, update)

	return err
}
