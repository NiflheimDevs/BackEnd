package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/services"
)

type FileHandler struct {
	FileService *services.FileService
	Validator   *validator.Validate
	JWTService  *services.JWT
	Constants   *bootstrap.Constants
}

func NewFileHandler(
	constants *bootstrap.Constants,
	fileService *services.FileService,
	jwtService *services.JWT,
	validator *validator.Validate,
) *FileHandler {
	return &FileHandler{
		Constants:   constants,
		FileService: fileService,
		Validator:   validator,
		JWTService:  jwtService,
	}
}

func (fh *FileHandler) UploadProfilePhoto(w http.ResponseWriter, r *http.Request) {
	type Photo struct {
		Data []byte `json:"data" validate:"required"`
	}
	userid := r.Context().Value(fh.Constants.Context.UserID).(int)
	params := Validated[Photo](fh.Validator, r)

}
