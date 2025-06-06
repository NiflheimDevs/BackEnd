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

type TeamElastic struct {
	Es *elasticsearch.Client
}

func NewTeamElastic(es *elasticsearch.Client) *TeamElastic {
	return &TeamElastic{
		Es: es,
	}
}

func (pe *TeamElastic) UpsertTeamDoc(teamem *elasticmodel.TeamElasticModel) error {
	upsertBody := map[string]interface{}{
		"doc":           teamem,
		"doc_as_upsert": true,
	}

	body, err := json.Marshal(upsertBody)
	if err != nil {
		log.Println("Error: failed to parse upsert request to json. details:", err)
		return err
	}
	docid := strconv.Itoa(teamem.ID)

	res, err := pe.Es.Update(enums.ELASTIC_TEAM_INDEX, docid, bytes.NewReader(body))
	if err != nil {
		log.Println("UpsertError: failed to upsert team document (from go). details:", err, "\nthe request:", body)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {

		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("UpsertError: failed to upsert team document (from elastic). response:", e, "\nthe request:", body)
		return errors.New("elastic error")
	}
	return nil
}

func (pe *TeamElastic) GetTeamDoc(teamid int) *elasticmodel.TeamElasticModel {
	var pem elasticmodel.ElasticsearchHit[elasticmodel.TeamElasticModel]

	res, err := pe.Es.Get(enums.ELASTIC_TEAM_INDEX, strconv.Itoa(teamid))
	if err != nil {
		log.Println("GetError: go error: failed to get team. detail:", err)
		return nil
	}
	defer res.Body.Close()
	if res.IsError() {

		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("GetError: elastic error: failed to get team document. response:", e)
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(&pem)
	if err != nil {
		log.Println("ParseError: get team. details: ", err)
	}
	return &pem.Source
}

func (pe *TeamElastic) DeleteTeamDoc(teamid int) error {
	res, err := pe.Es.Delete(enums.ELASTIC_TEAM_INDEX, strconv.Itoa(teamid))
	if err != nil {
		log.Println("DeleteError: go error: failed to delete team. detail:", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("DeleteError: elastic error: failed to delete team document. response:", e)
		return errors.New("elastic error")
	}
	return nil
}
