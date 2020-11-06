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

type leaderSelectionRequest struct {
	Member string `json:"member"`
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

func makePlayerActionRequest(url string, body interface{}, session string) {
	var bodyJson []byte
	var err error
	if body != nil {
		bodyJson, err = json.Marshal(body)
		if err != nil {
			log.Fatalf("can't marshall body: %+v\n", err)
		}
	}
	client := http.Client{}
	request, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:44324/%s", url), bytes.NewReader(bodyJson))
	request.AddCookie(&http.Cookie{Name: "session", Value: session})
	if err != nil {
		log.Fatalf("can't create request: %+v\n", err)
	}
	_, err = client.Do(request)
	if err != nil {
		log.Fatalf("can't do request: %+v\n", err)
	}
}

func StartGame(session string) {
	makePlayerActionRequest("actions/start-game", nil, session)
}

func LeaderSelectsMember(session, name string) {
	makePlayerActionRequest("actions/leader-selects-member", leaderSelectionRequest{Member: name}, session)
}

func LeaderDeselectsMember(session, name string) {
	makePlayerActionRequest("actions/leader-deselects-member", leaderSelectionRequest{Member: name}, session)
}

func LeaderConfirmsTeam(session string) {
	makePlayerActionRequest("actions/leader-confirms-team", nil, session)
}

func ApproveTeam(session string) {
	makePlayerActionRequest("actions/approve-team", nil, session)
}

func RejectTeam(session string) {
	makePlayerActionRequest("actions/reject-team", nil, session)
}

func SucceedMission(session string) {
	makePlayerActionRequest("actions/succeed-mission", nil, session)
}

func FailMission(session string) {
	makePlayerActionRequest("actions/fail-mission", nil, session)
}
