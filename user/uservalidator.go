package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/rosytucker/codenight/web"
)

func ValidUserId(objectId string) bool {
	return govalidator.IsMongoID(objectId)
}

func (current PublicUser) Validate() []web.ValidationError {
	errors := make([]web.ValidationError, 0)

	if len(current.AvatarUrl) > 0 && !govalidator.IsURL(current.AvatarUrl) {
		validationError := web.ValidationError{Field: "avatarUrl", Message: "invalid url supplied"}
		errors = append(errors, validationError)
	}
	return errors
}
