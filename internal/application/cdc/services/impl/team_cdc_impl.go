package cdcservicesimpl

import (
	"bytes"
	"encoding/json"
	"log"

	cdcdto "github.com/niflheimdevs/backend/internal/application/cdc/dto"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type TeamCdc struct {
	TagRepo repositories.TagRepo
}

func NewTeamCdc(tagRepo repositories.TagRepo) *TeamCdc {
	return &TeamCdc{
		TagRepo: tagRepo,
	}
}

func (tc *TeamCdc) TeamCapturer(data []byte) error {
	var message cdcdto.TeamCdcMessage
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
		log.Println("JsonDecodingError: Failed to Decode project message. detail:", err)
		return err
	}

	if message.Type == "d" {

	} else if message.Type == "c" || message.Type == "u" {

	}

	return nil
}
