package pkg

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

func NewValidator() *validator.Validate {
	MyValidator := validator.New()

	MyValidator.RegisterValidation("password", passwordValidation)
	MyValidator.RegisterValidation("username", usernameValidation)
	MyValidator.RegisterValidation("phone", phoneValidation)

	return MyValidator
}

func validateRegex(errors *exceptions.Exception, regex string, text string, tag exceptions.SpecificError) {
	matched, _ := regexp.MatchString(regex, text)
	if !matched {
		errors.AddError(tag)
	}
}

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var exc exceptions.Exception
	exc.Tag = exceptions.BAD_REQUEST
	validateRegex(&exc, "^.{0,32}$", password, exceptions.PASSWORD_TOO_LONG)
	validateRegex(&exc, "^.{8,}$", password, exceptions.PASSWORD_TOO_SHORT)
	validateRegex(&exc, "[a-z]", password, exceptions.PASSWORD_NO_SMALL)
	validateRegex(&exc, "[A-Z]", password, exceptions.PASSWORD_NO_CAPITAL)
	validateRegex(&exc, "[0-9]", password, exceptions.PASSWORD_NO_DIGIT)
	if len(exc.Errors) > 0 {
		panic(exc)
	}
	return true
}

func usernameValidation(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	var exc exceptions.Exception
	exc.Tag = exceptions.BAD_REQUEST

	validateRegex(&exc, "^.{0,32}$", username, exceptions.USERNAME_TOO_LONG)
	validateRegex(&exc, "^.{2,}$", username, exceptions.USERNAME_TOO_SHORT)

	if len(exc.Errors) > 0 {
		panic(exc)
	}
	return true
}

func phoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	var exc exceptions.Exception

	validateRegex(&exc, "^09\\d{9}$", phone, exceptions.PHONE_INVALID)
	if len(exc.Errors) > 0 {
		panic(exc)
	}
	return true
}
