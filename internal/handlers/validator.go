package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validated[T any](r *http.Request) T {
	var params T
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		panic(err)
	}

	log.Println(params)

	if err := validate.Struct(params); err != nil {
		panic(err)
	}

	return params
}
