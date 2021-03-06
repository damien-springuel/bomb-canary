package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/damien-springuel/bomb-canary/server/clientstream"
	"github.com/damien-springuel/bomb-canary/server/codegenerator"
	"github.com/damien-springuel/bomb-canary/server/gamehub"
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
	"github.com/damien-springuel/bomb-canary/server/messagelogger"
	"github.com/damien-springuel/bomb-canary/server/party"
	"github.com/damien-springuel/bomb-canary/server/playeractions"
	"github.com/damien-springuel/bomb-canary/server/sessions"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gookit/color"
)

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

var blackOnGreen = color.Style{color.BgLightGreen, color.FgBlack}
var blackOnBlue = color.Style{color.BgLightBlue, color.FgBlack}
var blackOnYellow = color.Style{color.BgLightYellow, color.FgBlack}

type colorPrinter struct{}

func (c colorPrinter) PrintCommand(m messagebus.Message) {
	blackOnGreen.Printf("%s - Command: %#v\n", time.Now().Format("2006/01/02 - 15:04:05"), m)
}
func (c colorPrinter) PrintEvent(m messagebus.Message) {
	blackOnBlue.Printf("%s - Event: %#v\n", time.Now().Format("2006/01/02 - 15:04:05"), m)
}

func main() {
	bus := messagebus.NewMessageBus()
	defer bus.Close()

	bus.SubscribeConsumer(messagelogger.New(colorPrinter{}))

	nowRandom := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomRuneGenerator := func() rune {
		min := 65 // A
		max := 90 // Z
		return rune(nowRandom.Intn(max-min+1) + min)
	}
	hub := gamehub.New(codegenerator.New(randomRuneGenerator), bus, randomAllegianceGenerator{})
	bus.SubscribeConsumer(hub)

	// sessions := sessions.New(uuidV4{})
	sessions := sessions.New(&easySession{}) // for easy testing purposes

	clientStreamer := clientstream.NewClientsStreamer(bus)
	eventReplayer := clientstream.NewEventReplayer(clientStreamer)
	clientEventBroker := clientstream.NewClientEventBroker(eventReplayer)
	bus.SubscribeConsumer(clientEventBroker)
	bus.SubscribeConsumer(eventReplayer)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = []string{"http://localhost:44322", "http://192.168.0.103:44322"}
	router.Use(cors.New(corsConfig))

	party.Register(router, party.NewPartyService(hub, hub, bus), sessions)
	playeractions.Register(router, sessions, playeractions.NewActionService(bus))
	clientstream.Register(router, sessions, clientStreamer)

	router.StaticFile("/", "./index.html")

	port := ":44324"
	blackOnYellow.Printf("serving %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		blackOnYellow.Printf("error serving %v\n", err)
	}
}
