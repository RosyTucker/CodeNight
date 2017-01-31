package user

import (
	"gopkg.in/mgo.v2/bson"
)

func ValidateId(objectId string) bool {
	return bson.IsObjectIdHex(objectId)
}
