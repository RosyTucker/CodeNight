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

	if !govalidator.IsByteLength(current.AvatarUrl, 0, 600) {
		validationError := web.ValidationError{Field: "avatarUrl", Message: "invalid length, max 600"}
		errors = append(errors, validationError)
	}

	if len(current.AvatarUrl) > 0 && !govalidator.IsURL(current.AvatarUrl) {
		validationError := web.ValidationError{Field: "avatarUrl", Message: "invalid url supplied"}
		errors = append(errors, validationError)
	}

	if !govalidator.IsByteLength(current.Blog, 0, 600) {
		validationError := web.ValidationError{Field: "blog", Message: "invalid length, max 600"}
		errors = append(errors, validationError)
	}

	if len(current.Blog) > 0 && !govalidator.IsURL(current.Blog) {
		validationError := web.ValidationError{Field: "blog", Message: "invalid url supplied"}
		errors = append(errors, validationError)
	}

	if !govalidator.IsByteLength(current.Description, 0, 600) {
		validationError := web.ValidationError{Field: "description", Message: "invalid length, max 600"}
		errors = append(errors, validationError)
	}

	if !govalidator.IsByteLength(current.Location, 0, 255) {
		validationError := web.ValidationError{Field: "location", Message: "invalid length, max 255"}
		errors = append(errors, validationError)
	}
	if !govalidator.IsByteLength(current.Name, 0, 255) {
		validationError := web.ValidationError{Field: "name", Message: "invalid length, max 255"}
		errors = append(errors, validationError)
	}

	return errors
}
