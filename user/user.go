package user

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        *string       `json:"name" bson:"name"`
	Token       string        `json:"-" bson:"token"`
	UserName    string        `json:"username" bson:"username"`
	Email       *string       `json:"email" bson:"email"`
	Description *string       `json:"description" bson:"description"`
	Blog        *string       `json:"blog" bson:"blog"`
	Location    *string       `json:"location" bson:"location"`
	AvatarUrl   *string       `json:"avatar_url" bson:"avatar_url"`
	IsAdmin     bool          `json:"isAdmin" bson:"is_admin"`
}

type PublicUser struct {
	Id          bson.ObjectId
	Name        *string `json:"name"`
	UserName    string  `json:"username"`
	Description *string `json:"description"`
	Blog        *string `json:"blog"`
	Location    *string `json:"location"`
	AvatarUrl   *string `json:"avatar_url"`
}

func FromJsonBody(userBody io.ReadCloser) (*User, error) {
	defer userBody.Close()
	decoder := json.NewDecoder(userBody)
	var user *User
	err := decoder.Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
