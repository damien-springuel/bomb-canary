package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type createPartyRequest struct {
	Name string `json:"name"`
}

type createPartyResponse struct {
	Code string `json:"code"`
}

func main() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	l := widgets.NewList()
	l.Title = "Actions"
	l.Rows = []string{
		"[0] Create Game",
	}
	l.TextStyle = termui.NewStyle(termui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)

	termui.Render(l)

	uiEvents := termui.PollEvents()
	go func() {
		for {
			e := <-uiEvents
			switch e.ID {
			case "0":
				log.Printf("another channel reader \n")
			}
		}
	}()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "0":
			code, session := createGame("From Cli")
			log.Printf("code: %s session: %s\n", code, session)
		}

		termui.Render(l)
	}
}

func createGame(name string) (code string, session string) {
	createPartyJson, err := json.Marshal(createPartyRequest{Name: "From CLI"})
	if err != nil {
		log.Fatalf("can't marshall create party request: %+v\n", err)
	}
	client := http.Client{}
	request, err := http.NewRequest("POST", "http://localhost:44324/party/create", bytes.NewReader(createPartyJson))
	if err != nil {
		log.Fatalf("can't create request: %+v\n", err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("can't do request: %+v\n", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("can't read response body: %+v\n", err)
	}
	actualResponse := createPartyResponse{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		log.Fatalf("can't unmarshall response: %+v\n", err)
	}
	code = actualResponse.Code
	cookies := response.Cookies()
	for _, c := range cookies {
		if c.Name == "session" {
			session = c.Value
			break
		}
	}
	return
}
