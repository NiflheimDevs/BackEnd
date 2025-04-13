package servicesimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
)

type SmsService struct {
	Env  *bootstrap.Env
	Cons *bootstrap.Constants
}

func NewSmsService(env *bootstrap.Env) *SmsService {
	return &SmsService{Env: env}
}

func (ss *SmsService) GenerateOTP() string {
	return fmt.Sprintf("%05d", rand.Int()%100000)

}

func (ss *SmsService) SendOTP(phonenumber string, code string) {
	var apikey string

	if ss.Cons.DevelopMode {
		apikey = ss.Env.SMS.APISandboxKey
	} else {
		apikey = ss.Env.SMS.APIKey
	}

	type Parameters struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	type Template struct {
		Phonenumber string       `json:"mobile"`
		TemplateId  string       `json:"templateid"`
		Parameter   []Parameters `json:"parameters"`
	}
	temp := Template{
		Phonenumber: phonenumber,
		TemplateId:  "485976",
		Parameter: []Parameters{Parameters{
			Name:  "Code",
			Value: code,
		}},
	}
	marshalled, _ := json.Marshal(temp)
	bodyReader := bytes.NewReader(marshalled)
	req, _ := http.NewRequest(http.MethodPost, ss.Env.SMS.IPAddr, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("x-api-key", apikey)

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
	data, _ := io.ReadAll(res.Body)
	log.Println(res.Status)
	log.Println(string(data))
}
