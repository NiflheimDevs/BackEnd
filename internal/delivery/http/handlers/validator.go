package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

func Validated[T any](validate *validator.Validate, r *http.Request) T {
	var params T
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	log.Println(params)

	if err := validate.Struct(params); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.BAD_REQUEST,
			Errors: []exceptions.SpecificError{exceptions.MISSING_REQUIRED_FIELD},
		})

	}

	return params
}
