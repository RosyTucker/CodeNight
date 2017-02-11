package user

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Token       string        `json:"-" bson:"token"`
	UserName    string        `json:"username" bson:"username"`
	Email       string        `json:"email" bson:"email"`
	Description string        `json:"description" bson:"description"`
	Blog        string        `json:"blog" bson:"blog"`
	Location    string        `json:"location" bson:"location"`
	AvatarUrl   string        `json:"avatarUrl" bson:"avatar_url"`
	IsAdmin     bool          `json:"isAdmin" bson:"is_admin"`
}

type PublicUser struct {
	Id          bson.ObjectId
	Name        string `json:"name"`
	Description string `json:"description"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	AvatarUrl   string `json:"avatarUrl"`
}

func PublicFromJsonBody(userBody io.ReadCloser) (*PublicUser, error) {
	defer userBody.Close()
	decoder := json.NewDecoder(userBody)
	var public *PublicUser
	err := decoder.Decode(&public)
	if err != nil {
		return nil, err
	}

	return public, nil
}
