package servicesimpl

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type CommentService struct {
	CommentRepo    repositories.CommentRepo
	ProjectService services.ProjectService
}

func NewCommentService(
	commentRepo repositories.CommentRepo,
	projectService services.ProjectService,
) *CommentService {
	return &CommentService{
		CommentRepo:    commentRepo,
		ProjectService: projectService,
	}
}

func (cs CommentService) PutComment(userID, ProjectID int, content string, star int) int {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_TOKEN_EXPIRED},
		})
	}

	project := cs.ProjectService.GetProject(ProjectID)

	if project.State != 4 {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	if userID != project.OwnerID {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.USER_NOT_OWNER},
		})
	}

	id := cs.CommentRepo.AddComment(ProjectID, project.SelectedBid, content, star)

	return id
}

func (cs CommentService) GetUserStar(userID int) float32 {
	star := cs.CommentRepo.GetStar(userID)
	return star
}

func (cs CommentService) GetUserComments(userID int) []dto.CommentDTO {
	comments := cs.CommentRepo.GetUserComments(userID)
	return comments
}

func (cs CommentService) GetCommentInfo(id int) dto.CommentDTO {
	comment := cs.CommentRepo.GetCommentInfo(id)
	return comment
}
