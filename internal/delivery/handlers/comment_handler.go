package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
)

type CommentHandler struct {
	CommentService services.CommentService
	Validator      *validator.Validate
	Constants      *bootstrap.Constants
}

func NewCommentHandler(
	commentService services.CommentService,
	validator *validator.Validate,
	constants *bootstrap.Constants,
) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
		Validator:      validator,
		Constants:      constants,
	}
}

func (ch *CommentHandler) PutComment(w http.ResponseWriter, r *http.Request) {
	type commentInfo struct {
		ProjectID int    `json:"project_id" validate:"required"`
		Content   string `json:"content" validate:"required"`
		Star      int    `json:"rating" validate:"required,lt=6,gt=0"`
	}

	params := Validated[commentInfo](ch.Validator, r)

	userid := r.Context().Value(ch.Constants.Context.UserID).(int)

	commentid := ch.CommentService.PutComment(userid, params.ProjectID, params.Content, params.Star)

	id := dto.PutComment{
		ID: commentid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
}

func (ch *CommentHandler) GetStar(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "user_id")
	userID, _ := strconv.Atoi(userIDString)

	star := ch.CommentService.GetUserStar(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(star)
}

func (ch *CommentHandler) GetUserComments(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "user_id")
	userID, _ := strconv.Atoi(userIDString)

	comments := ch.CommentService.GetUserComments(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (ch *CommentHandler) GetCommentInfo(w http.ResponseWriter, r *http.Request) {
	commentIDString := chi.URLParam(r, "id")
	commentID, _ := strconv.Atoi(commentIDString)

	comment := ch.CommentService.GetCommentInfo(commentID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}
