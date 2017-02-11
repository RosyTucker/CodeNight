package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/rosytucker/codenight/web"
	"gopkg.in/mgo.v2/bson"
)

func ValidUserId(objectId string) bool {
	return bson.IsObjectIdHex(objectId)
}

func (current *User) Validate() []web.ValidationError {
	errors := make([]web.ValidationError, 0)

	if !govalidator.IsEmail(current.Email) {
		validationError := web.ValidationError{Field: "email", Message: "Invalid format"}
		errors = append(errors, validationError)
	}
	return errors
}
