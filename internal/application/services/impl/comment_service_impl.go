package servicesimpl

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type CommentService struct {
	CommentRepo    repositories.CommentRepo
	ProjectRepo    repositories.ProjectRepo
	ProjectService services.ProjectService
}

func NewCommentService(
	commentRepo repositories.CommentRepo,
	projectRepo repositories.ProjectRepo,
	projectService services.ProjectService,
) *CommentService {
	return &CommentService{
		CommentRepo:    commentRepo,
		ProjectRepo:    projectRepo,
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

	cs.ProjectRepo.UpdateProjectState(ProjectID, 5)

	return id
}

func (cs CommentService) GetUserStarAndComment(userID int) (float64, int) {
	star, err := cs.CommentRepo.GetStar(userID)
	commentCount, err := cs.CommentRepo.GetCommentCount(userID)
	if err != nil {
		return 0, 0
	}
	return star, commentCount
}

func (cs CommentService) GetUserComments(userID int) []dto.CommentWithUserDTO {
	comments, err := cs.CommentRepo.GetUserComments(userID)
	if err != nil {
		return nil
	}
	return comments
}

func (cs CommentService) GetCommentInfo(id int) *dto.CommentWithUserDTO {
	comment := cs.CommentRepo.GetCommentInfo(id)
	return &comment
}

func (cs CommentService) GetCommentOfProject(projectID int) *dto.CommentDTO {
	comment, err := cs.CommentRepo.GetCommentOfProject(projectID)

	if err != nil {
		return nil
	}

	dto := &dto.CommentDTO{
		ID:        comment.ID,
		ProjectID: comment.ProjectID,
		Content:   comment.Content,
		Rating:    comment.Rating,
	}

	return dto
}
