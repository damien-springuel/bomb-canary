package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/damien-springuel/bomb-canary/server/gamehub"
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
	"github.com/damien-springuel/bomb-canary/server/messagelogger"
	"github.com/damien-springuel/bomb-canary/server/party"
	"github.com/damien-springuel/bomb-canary/server/playeractions"
	"github.com/damien-springuel/bomb-canary/server/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gookit/color"
)

type randomCodeGenerator struct{}

func (r randomCodeGenerator) GenerateCode() string {
	return "AAA"
}

type randomAllegianceGenerator struct{}

func (r randomAllegianceGenerator) Generate(nbPlayers, nbSpies int) []gamerules.Allegiance {
	allegiances := make([]gamerules.Allegiance, nbPlayers)
	for i := range allegiances {
		if i < nbSpies {
			allegiances[i] = gamerules.Spy
		} else {
			allegiances[i] = gamerules.Resistance
		}
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(allegiances), func(i, j int) { allegiances[i], allegiances[j] = allegiances[j], allegiances[i] })

	return allegiances
}

type uuidV4 struct{}

func (u uuidV4) Create() string {
	return uuid.New().String()
}

type easySession struct {
	mut       *sync.Mutex
	currentId int
}

func (e *easySession) Create() string {
	if e.mut == nil {
		e.mut = &sync.Mutex{}
	}
	e.mut.Lock()
	defer e.mut.Unlock()
	e.currentId += 1
	session := fmt.Sprintf("%d", e.currentId)
	return session
}

type colorPrinter struct{}

func (c colorPrinter) PrintCommand(m messagebus.Message) {
	color.Style{color.BgLightGreen, color.FgBlack}.Printf("%s - Command: %#v\n", time.Now().Format("2006/01/02 - 15:04:05"), m)
}
func (c colorPrinter) PrintEvent(m messagebus.Message) {
	color.Style{color.BgLightBlue, color.FgBlack}.Printf("%s - Event: %#v\n", time.Now().Format("2006/01/02 - 15:04:05"), m)
}

func main() {
	bus := messagebus.NewMessageBus()

	hub := gamehub.New(randomCodeGenerator{}, bus, randomAllegianceGenerator{})
	bus.SubscribeConsumer(hub)
	bus.SubscribeConsumer(messagelogger.New(colorPrinter{}))

	// sessions := sessions.New(uuidV4{})
	sessions := sessions.New(&easySession{}) // for easy testing purposes

	router := gin.Default()
	party.Register(router, party.NewPartyService(hub, hub, bus), sessions)
	playeractions.Register(router, sessions, playeractions.NewActionService(bus))

	port := ":44324"
	log.Printf("serving %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("error serving %v\n", err)
	}
}
