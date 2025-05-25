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

type UserCdc struct {
	TagRepo     repositories.TagRepo
	UserElastic elastic.UserElastic
}

func NewUserCdc(tagRepo repositories.TagRepo, ue elastic.UserElastic) *UserCdc {
	return &UserCdc{
		TagRepo:     tagRepo,
		UserElastic: ue,
	}
}

func (uc *UserCdc) UserCapturer(data []byte) error {
	var message cdcdto.UserCapturer
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
		log.Println("JsonDecodingError: Failed to Decode user message. detail:", err)
		return err
	}

	if message.Type == "d" {

	} else if message.Type == "c" || message.Type == "u" {

	}

	return nil
}

func (uc *UserCdc) UserTagCapturer(data []byte) error {
	var message cdcdto.UserCareerTagCapturer
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
		log.Println("JsonDecodingError: Failed to Decode user message. detail:", err)
		return err
	}

	var userid, tagid int
	switch message.Type {
	case "d":
		if message.Before.Type != 0 {
			return nil
		}
		tagid = message.Before.Tagid
		userid = message.Before.UserCareerid
	case "c":
		if message.After.Type != 0 {
			return nil
		}
		tagid = message.After.Tagid
		userid = message.After.UserCareerid
	case "u":
		if message.After.Type != 0 {
			return nil
		}
		tagid = message.After.Tagid
		userid = message.After.UserCareerid
	default:
		return errors.New("not supported operation")
	}

	user := uc.UserElastic.GetUserDoc(userid)
	if user == nil {
		return errors.New("user is nil")
	}

	tag, err := uc.TagRepo.GetTag(tagid)
	if err != nil {
		return err
	}

	if message.Type == "d" {
		for i := 0; i < len(user.Tags); i++ {
			if user.Tags[i].Name == tag.Name {
				user.Tags = utils.RemoveUnordered(user.Tags, &i)
				break
			}
		}
	} else if message.Type == "c" {
		user.Tags = append(user.Tags, elasticmodel.TagElasticModel{
			Name:  tag.Name,
			Level: message.After.Level,
		})
	} else if message.Type == "u" {
		for i := 0; i < len(user.Tags); i++ {
			if user.Tags[i].Name == tag.Name {
				user.Tags[i].Level = message.After.Level
				break
			}
		}
	}

	return uc.UserElastic.UpsertUserDoc(user)

}
