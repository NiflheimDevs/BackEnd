package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
)

type FileHandler struct {
	FileService services.FileService
	TeamService services.TeamService
	Validator   *validator.Validate
	JWTService  services.JWT
	Constants   *bootstrap.Constants
	ENV         *bootstrap.Env
}

func NewFileHandler(
	constants *bootstrap.Constants,
	Env *bootstrap.Env,
	fileService services.FileService,
	teamService services.TeamService,
	jwtService services.JWT,
	validator *validator.Validate,
) *FileHandler {
	return &FileHandler{
		Constants:   constants,
		ENV:         Env,
		FileService: fileService,
		Validator:   validator,
		JWTService:  jwtService,
		TeamService: teamService,
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

	fh.FileService.UploadProfilePhoto(data, userid)

	type PhotoResponse struct {
		PhotoPath string `json:"photo_path"`
	}

	var photoPath PhotoResponse

	photoPath.PhotoPath = fh.FileService.GetProfilePhotoURL(userid, true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photoPath)
}

func (fh *FileHandler) DeleteProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(fh.Constants.Context.UserID).(int)
	fh.FileService.DeleteProfilePhoto(userid)
	w.WriteHeader(http.StatusNoContent)
}

func (fh *FileHandler) UploadTeamProfilePhoto(w http.ResponseWriter, r *http.Request) {
	teamidString := chi.URLParam(r, "team_id")
	teamid, err := strconv.Atoi(teamidString)
	if err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

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

	fh.TeamService.UpdateTeamProfile(userid, int64(teamid), data)

	type PhotoResponse struct {
		PhotoPath string `json:"photo_path"`
	}

	var photoPath PhotoResponse

	photoPath.PhotoPath = fh.FileService.GetTeamProfilePhotoURL(int64(teamid), true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photoPath)
}

func (fh *FileHandler) DeleteTeamProfilePhoto(w http.ResponseWriter, r *http.Request) {
	teamidString := chi.URLParam(r, "team_id")
	teamid, err := strconv.Atoi(teamidString)
	if err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	userid := r.Context().Value(fh.Constants.Context.UserID).(int)

	fh.TeamService.DeleteTeamProfile(userid, int64(teamid))

	w.WriteHeader(http.StatusNoContent)

}

// func (fh *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
// 	// automatically generates 404 if file not found (based on the document)
// 	http.StripPrefix("/storage/", http.FileServer(http.Dir(fh.Constants.StorageDir))).ServeHTTP(w, r)
// }

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

func (fh *FileHandler) GetProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(fh.Constants.Context.UserID).(int)
	high_qual := fh.FileService.GetProfilePhotoURL(userid, true)
	low_qual := fh.FileService.GetProfilePhotoURL(userid, false)
	type profilepic struct {
		HighQuality string `json:"high_quality"`
		LowQuality  string `json:"low_quality"`
	}

	pic := profilepic{
		HighQuality: high_qual,
		LowQuality:  low_qual,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pic); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

// func (fh *FileHandler) GetTeamProfilePhoto(w http.ResponseWriter, r *http.Request) {

// 	high_qual := fh.FileService.GetProfilePhotoURL(userid, true)
// 	low_qual := fh.FileService.GetProfilePhotoURL(userid, false)
// 	type profilepic struct {
// 		HighQuality string `json:"high_quality"`
// 		LowQuality  string `json:"low_quality"`
// 	}

// 	pic := profilepic{
// 		HighQuality: high_qual,
// 		LowQuality:  low_qual,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(pic); err != nil {
// 		panic(exceptions.Exception{
// 			Tag:    exceptions.INTERNAL_ERROR,
// 			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
// 		})
// 	}
// }

func (fh *FileHandler) GetResume(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(fh.Constants.Context.UserID).(int)

	url := fh.FileService.GetResumeURL(userid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url))
}
