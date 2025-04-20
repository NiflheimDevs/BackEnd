package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
)

type FileHandler struct {
	FileService services.FileService
	Validator   *validator.Validate
	JWTService  services.JWT
	Constants   *bootstrap.Constants
	ENV         *bootstrap.Env
}

func NewFileHandler(
	constants *bootstrap.Constants,
	Env *bootstrap.Env,
	fileService services.FileService,
	jwtService services.JWT,
	validator *validator.Validate,
) *FileHandler {
	return &FileHandler{
		Constants:   constants,
		ENV:         Env,
		FileService: fileService,
		Validator:   validator,
		JWTService:  jwtService,
	}
}

func (fh *FileHandler) UploadProfilePhoto(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	defer file.Close()

	if header.Size > fh.Constants.MaxPhotoSize {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.FILE_TOO_LARGE},
		})
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	userid := r.Context().Value(fh.Constants.Context.UserID).(int)

	res := fh.FileService.UploadProfilePhoto(data, userid)

	type PhotoResponse struct {
		PhotoPath string `json:"photo_path"`
	}

	var photoPath PhotoResponse

	photoPath.PhotoPath = fmt.Sprintf("%s/storage/%s", fh.ENV.Server.IP_Addr, res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photoPath)
}

func (fh *FileHandler) DeleteProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(fh.Constants.Context.UserID).(int)
	fh.FileService.DeleteProfilePhoto(userid)
	w.WriteHeader(http.StatusNoContent)
}

func (fh *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	// automatically generates 404 if file not found (based on the document)
	http.StripPrefix("/storage/", http.FileServer(http.Dir(fh.Constants.StorageDir))).ServeHTTP(w, r)
}

func (fh *FileHandler) UploadUserResume(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	defer file.Close()

	if header.Size > fh.Constants.MaxResumeSize {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.FILE_TOO_LARGE},
		})
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	userid := r.Context().Value(fh.Constants.Context.UserID).(int)

	fh.FileService.UploadResume(data, userid)

	w.WriteHeader(http.StatusNoContent)
}

func (fh *FileHandler) DeleteUserResume(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(fh.Constants.Context.UserID).(int)
	fh.FileService.DeleteResume(userid)
	w.WriteHeader(http.StatusNoContent)
}
