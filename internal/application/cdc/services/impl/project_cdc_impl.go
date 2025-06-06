package cdcservicesimpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	cdcdto "github.com/niflheimdevs/backend/internal/application/cdc/dto"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
	"github.com/niflheimdevs/backend/internal/domain/repositories/elastic"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/utils"
)

type ProjectCdc struct {
	TagRepo        repositories.TagRepo
	ProjectElastic elastic.ProjectElastic
}

func NewProjectCdc(tagRepo repositories.TagRepo, pe elastic.ProjectElastic) *ProjectCdc {
	return &ProjectCdc{
		TagRepo:        tagRepo,
		ProjectElastic: pe,
	}
}

func (pc *ProjectCdc) ProjectCapturer(data []byte) error {
	var message cdcdto.ChangeDataCaptureDto[elasticmodel.ProjectElasticModel]
	var err error
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
		log.Println("JsonDecodingError: Failed to Decode project message. detail:", err)
		return err
	}

	if message.Type == "d" {
		err = pc.ProjectElastic.DeleteProjectDoc(message.Before.ID)

	} else if message.Type == "c" || message.Type == "u" {
		err = pc.ProjectElastic.UpsertProjectDoc(message.After)
	}

	return err

}

func (pc *ProjectCdc) ProjectTagCapturer(data []byte) error {
	var message cdcdto.ChangeDataCaptureDto[cdcdto.ProjectTagDto]
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
		log.Println("JsonDecodingError: Failed to Decode project message. detail:", err)
		return err
	}

	var projectid, tagid int
	switch message.Type {
	case "d":
		projectid = message.Before.Projectid
		tagid = message.Before.Tagid
	case "c":
		tagid = message.After.Tagid
		projectid = message.After.Projectid
	default:
		return errors.New("not supported operation (u is unsupported)")
	}

	project := pc.ProjectElastic.GetProjectDoc(projectid)
	if project == nil {
		return errors.New("project is nil")
	}

	tag, err := pc.TagRepo.GetTag(tagid)
	if err != nil {
		return err
	}

	//if multiple instances of this code start to work, then you should fetch tags from the main database in order to remain up to date
	if message.Type == "d" {
		for i := 0; i < len(project.Tags); i++ {
			if project.Tags[i] == tag.Name {
				project.Tags = utils.RemoveUnordered[string](project.Tags, &i)
				break
			}
		}
	} else if message.Type == "c" {
		project.Tags = append(project.Tags, tag.Name)
	}

	return pc.ProjectElastic.UpsertProjectDoc(project)
}
