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
	l.Title = "Actions (^C to quit)"
	l.Rows = []string{
		"[0] Create Game",
	}

	grid := termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(termui.NewRow(1, l))

	termui.Render(grid)

	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(termui.Resize)
			grid.SetRect(0, 0, payload.Width, payload.Height)
			termui.Clear()
		case "0":
			code, session := createGame("From Cli")
			log.Printf("code: %s session: %s\n", code, session)
		}

		termui.Render(grid)
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
