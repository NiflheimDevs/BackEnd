package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

func NewValidator() *validator.Validate {
	MyValidator := validator.New()

	MyValidator.RegisterValidation("password", passwordValidation)
	MyValidator.RegisterValidation("username", usernameValidation)
	MyValidator.RegisterValidation("phone", phoneValidation)

	return MyValidator
}

func Validated[T any](validate *validator.Validate, r *http.Request) T {
	var params T
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: enums.BAD_REQUEST,
		})
	}

	log.Println(params)

	if err := validate.Struct(params); err != nil {
		panic(exceptions.Exception{
			Tag:    enums.BAD_REQUEST,
			Errors: []enums.SpecificError{enums.MISSING_REQUIRED_FIELD},
		})

	}

	return params
}

func validateRegex(errors *exceptions.Exception, regex string, text string, tag enums.SpecificError) {
	matched, _ := regexp.MatchString(regex, text)
	if !matched {
		errors.AddError(tag)
	}
}

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var exc exceptions.Exception
	exc.Tag = enums.BAD_REQUEST
	validateRegex(&exc, "^.{0,32}$", password, enums.PASSWORD_TOO_LONG)
	validateRegex(&exc, "^.{8,}$", password, enums.PASSWORD_TOO_SHORT)
	validateRegex(&exc, "[a-z]", password, enums.PASSWORD_NO_SMALL)
	validateRegex(&exc, "[A-Z]", password, enums.PASSWORD_NO_CAPITAL)
	validateRegex(&exc, "[0-9]", password, enums.PASSWORD_NO_DIGIT)
	if len(exc.Errors) > 0 {
		panic(exc)
	}
	return true
}

func usernameValidation(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	var exc exceptions.Exception
	exc.Tag = enums.BAD_REQUEST

	validateRegex(&exc, "^.{0,32}$", username, enums.USERNAME_TOO_LONG)
	validateRegex(&exc, "^.{2,}$", username, enums.USERNAME_TOO_SHORT)

	if len(exc.Errors) > 0 {
		panic(exc)
	}
	return true
}

func phoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	var exc exceptions.Exception

	validateRegex(&exc, "^09\\d{9}$", phone, enums.PHONE_INVALID)
	if len(exc.Errors) > 0 {
		panic(exc)
	}
	return true
}
