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

	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

func GenerateOTP() string {
	return fmt.Sprintf("%05d", rand.Int()%100000)

}

func SendOTP(phonenumber string, code string) {
	// apikey := "OQIAPP4fRTpqWpWafX2lljoW9YBSuCmGLdFGFDZfJCfLfc97"
	apikey := "i9jivkYg8ONebnmtTb5ncBcOuaFoCIxsUyyTWKVcOSXaK3da"
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
	req, _ := http.NewRequest(http.MethodPost, "https://api.sms.ir/v1/send/verify", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("x-api-key", apikey)

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
		})
	}
	data, _ := io.ReadAll(res.Body)
	log.Println(res.Status)
	log.Println(string(data))
}
