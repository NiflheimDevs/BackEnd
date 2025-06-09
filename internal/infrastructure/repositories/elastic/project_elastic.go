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

type ProjectElastic struct {
	Es *elasticsearch.Client
}

func NewProjectElastic(es *elasticsearch.Client) *ProjectElastic {
	return &ProjectElastic{
		Es: es,
	}
}

func (pe *ProjectElastic) UpsertProjectDoc(projectem *elasticmodel.ProjectElasticModel) error {
	upsertBody := map[string]interface{}{
		"doc":           projectem,
		"doc_as_upsert": true,
	}

	body, err := json.Marshal(upsertBody)
	if err != nil {
		log.Println("Error: failed to parse upsert request to json. details:", err)
		return err
	}

	docid := strconv.Itoa(projectem.ID)

	res, err := pe.Es.Update(enums.ELASTIC_PROJECT_INDEX, docid, bytes.NewReader(body))
	if err != nil {
		log.Println("UpsertError: failed to upsert project document (from go). details:", err, "\nthe request:", body)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {

		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("UpsertError: failed to upsert project document (from elastic). response:", e, "\nthe request:", body)
		return errors.New("elastic error")
	}
	return nil
}

func (pe *ProjectElastic) GetProjectDoc(projectid int) *elasticmodel.ProjectElasticModel {
	var pem elasticmodel.ElasticsearchHit[elasticmodel.ProjectElasticModel]

	res, err := pe.Es.Get(enums.ELASTIC_PROJECT_INDEX, strconv.Itoa(projectid))
	if err != nil {
		log.Println("GetError: go error: failed to get project. detail:", err)
		return nil
	}
	defer res.Body.Close()
	if res.IsError() {

		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("GetError: elastic error: failed to get project document. response:", e)
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(&pem)
	if err != nil {
		log.Println("ParseError: get project. details: ", err)
	}
	return &pem.Source
}

func (pe *ProjectElastic) DeleteProjectDoc(projectid int) error {
	res, err := pe.Es.Delete(enums.ELASTIC_PROJECT_INDEX, strconv.Itoa(projectid))
	if err != nil {
		log.Println("DeleteError: go error: failed to delete project. detail:", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		json.NewDecoder(res.Body).Decode(&e)
		log.Println("DeleteError: elastic error: failed to delete project document. response:", e)
		return errors.New("elastic error")
	}
	return nil
}
