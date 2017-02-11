package user

import (
	"gopkg.in/mgo.v2/bson"
)

func ValidUserId(objectId string) bool {
	return bson.IsObjectIdHex(objectId)
}
