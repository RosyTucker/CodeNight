package user

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type User struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Token       string        `json:"-" bson:"token"`
	UserName    string        `json:"username" bson:"username"`
	Email       string        `json:"email" bson:"email"`
	Description string        `json:"description" bson:"description"`
	Blog        string        `json:"blog" bson:"blog"`
	Location    string        `json:"location" bson:"location"`
	Company     string        `json:"company" bson:"company"`
	AvatarUrl   string        `json:"avatarUrl" bson:"avatar_url"`
	MemberSince string        `json:"memberSince" bson:"memberSince"`
	IsAdmin     bool          `json:"isAdmin" bson:"is_admin"`
}

type PublicUser struct {
	Id          bson.ObjectId `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Blog        string        `json:"blog"`
	Location    string        `json:"location"`
	Company     string        `json:"company"`
	AvatarUrl   string        `json:"avatarUrl"`
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
