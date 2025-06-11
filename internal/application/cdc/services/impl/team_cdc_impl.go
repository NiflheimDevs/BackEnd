package cdcservicesimpl

import (
	"bytes"
	"encoding/json"
	"log"

	cdcdto "github.com/niflheimdevs/backend/internal/application/cdc/dto"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
	"github.com/niflheimdevs/backend/internal/domain/repositories/elastic"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type TeamCdc struct {
	TagRepo     repositories.TagRepo
	TeamElastic elastic.TeamElastic
}

func NewTeamCdc(tagRepo repositories.TagRepo, te elastic.TeamElastic) *TeamCdc {
	return &TeamCdc{
		TagRepo:     tagRepo,
		TeamElastic: te,
	}
}

func (tc *TeamCdc) TeamCapturer(data []byte) error {
	var message cdcdto.ChangeDataCaptureDto[elasticmodel.TeamElasticModel]
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
		log.Println("JsonDecodingError: Failed to Decode project message. detail:", err)
		return err
	}

	var err error
	if message.Type == "d" {
		if message.Before.Type == 1 {
			return nil
		}
		err = tc.TeamElastic.DeleteTeamDoc(message.Before.ID)
	} else if message.Type == "c" || message.Type == "u" {
		if message.After.Type == 1 {
			return nil
		}
		err = tc.TeamElastic.UpsertTeamDoc(message.After)
	}

	return err
}
