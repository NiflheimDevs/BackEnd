package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
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
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.BAD_REQUEST,
		})
	}

	defer file.Close()

	if header.Size > 50000000 {
		panic(exceptions.Exception{
			Tag:    enums.UNPROCESSABLE,
			Errors: []enums.SpecificError{enums.FILE_TOO_LARGE},
		})
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.BAD_REQUEST,
		})
	}

	userid := r.Context().Value(fh.Constants.Context.UserID).(int)

	res := fh.FileService.UploadProfilePhoto(data, userid)

	type PhotoResponse struct {
		PhotoPath string `json:"photo_path"`
	}

	var photoPath PhotoResponse

	photoPath.PhotoPath = fmt.Sprintf("%s/storage/%s", fh.Constants.IPAddr, res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photoPath)
}

func (fh *FileHandler) DeleteProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(fh.Constants.Context.UserID).(int)
	fh.FileService.DeleteProfilePhoto(userid)
	w.WriteHeader(http.StatusOK)
}

func (fh *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	// automatically generates 404 if file not found (based on the document)
	http.StripPrefix("/storage/", http.FileServer(http.Dir(fh.Constants.StorageDir))).ServeHTTP(w, r)

}
