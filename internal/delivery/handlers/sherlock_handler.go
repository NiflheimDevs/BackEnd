package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
)

type SherlockHandler struct {
	SherlockService services.SherlockService
	Validator       *validator.Validate
}

func NewSherlockHandler(ss services.SherlockService, v *validator.Validate) *SherlockHandler {
	return &SherlockHandler{
		SherlockService: ss,
		Validator:       v,
	}
}

func (sh *SherlockHandler) Search(w http.ResponseWriter, r *http.Request) {

	params := Validated[elasticmodel.SearchRequest](sh.Validator, r)

	if params.Order != "" && params.Order != "desc" && params.Order != "asc" {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	res := sh.SherlockService.SearchEverything(&params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
