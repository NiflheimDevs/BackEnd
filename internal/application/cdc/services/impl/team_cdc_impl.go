package cdcservicesimpl

import (
	"bytes"
	"encoding/json"
	"log"

	cdcdto "github.com/niflheimdevs/backend/internal/application/cdc/dto"
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
	var message cdcdto.TeamCdcMessage
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
		log.Println("JsonDecodingError: Failed to Decode project message. detail:", err)
		return err
	}

	var err error
	if message.Type == "d" {
		err = tc.TeamElastic.DeleteTeamDoc(message.Before.ID)
	} else if message.Type == "c" || message.Type == "u" {
		err = tc.TeamElastic.UpsertTeamDoc(message.After)
	}

	return err
}
