package elasticimpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
)

type UserElastic struct {
	Es *elasticsearch.Client
}

func NewUserElastic(es *elasticsearch.Client) *UserElastic {
	return &UserElastic{
		Es: es,
	}
}

func (pe *UserElastic) UpsertUserDoc(userem *elasticmodel.UserElasticModel) error {
	upsertBody := map[string]interface{}{
		"doc":           userem,
		"doc_as_upsert": true,
	}

	body, err := json.Marshal(upsertBody)
	if err != nil {
		log.Println("Error: failed to parse upsert request to json. details:", err)
		return err
	}

	docid := strconv.Itoa(userem.ID)

	res, err := pe.Es.Update(enums.ELASTIC_USER_INDEX, docid, bytes.NewReader(body))
	if err != nil {
		log.Println("UpsertError: failed to upsert user document (from go). details:", err, "\nthe request:", body)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {

		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("UpsertError: failed to upsert user document (from elastic). response:", e, "\nthe request:", body)
		return errors.New("elastic error")
	}

	log.Println("respond\n", res.Body)
	return nil
}

func (pe *UserElastic) GetUserDoc(userid int) *elasticmodel.UserElasticModel {
	var pem elasticmodel.UserElasticModel

	res, err := pe.Es.Get(enums.ELASTIC_USER_INDEX, strconv.Itoa(userid))
	if err != nil {
		log.Println("GetError: go error: failed to get user. detail:", err)
		return nil
	}
	defer res.Body.Close()
	if res.IsError() {

		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("GetError: elastic error: failed to get user document. response:", e)
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(&pem)
	if err != nil {
		log.Println("ParseError: get user. details: ", err)
	}
	return &pem
}

func (pe *UserElastic) DeleteUserDoc(userid int) error {
	res, err := pe.Es.Delete(enums.ELASTIC_USER_INDEX, strconv.Itoa(userid))
	if err != nil {
		log.Println("DeleteError: go error: failed to delete user. detail:", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("DeleteError: elastic error: failed to delete user document. response:", e)
		return errors.New("elastic error")
	}
	return nil
}
