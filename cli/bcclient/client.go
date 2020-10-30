package bcclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type createPartyRequest struct {
	Name string `json:"name"`
}

type createPartyResponse struct {
	Code string `json:"code"`
}

type joinPartyRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func makePartyRequest(path string, body interface{}, responseValue interface{}) (session string) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("can't marshall body: %+v\n", err)
	}
	client := http.Client{}
	request, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:44324/%s", path), bytes.NewReader(bodyJson))
	if err != nil {
		log.Fatalf("can't create request: %+v\n", err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("can't do request: %+v\n", err)
	}
	if responseValue != nil {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("can't read response body: %+v\n", err)
		}
		err = json.Unmarshal(responseBody, responseValue)
		if err != nil {
			log.Fatalf("can't unmarshall response: %+v\n", err)
		}
	}
	cookies := response.Cookies()
	for _, c := range cookies {
		if c.Name == "session" {
			return c.Value
		}
	}
	panic("should have session")
}

func CreateGame(name string) (code string, session string) {
	actualResponse := createPartyResponse{}
	session = makePartyRequest("party/create", createPartyRequest{Name: name}, &actualResponse)
	code = actualResponse.Code
	return
}

func JoinGame(code, name string) (session string) {
	session = makePartyRequest("party/join", joinPartyRequest{Code: code, Name: name}, nil)
	return
}
