package user

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        *string       `json:"name", bson:"name"`
	Token       string        `json:"-" bson:"token"`
	UserName    *string       `json:"username", bson:"username"`
	Email       *string       `json:"email", bson:"email"`
	Description *string       `json:"description", bson:"description"`
	Blog        *string       `json:"blog", bson:"blog"`
	Location    *string       `json:"location", bson:"location"`
	AvatarUrl   *string       `json:"avatar", bson: "avatar"`
}
