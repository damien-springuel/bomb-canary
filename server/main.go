package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/damien-springuel/bomb-canary/server/gamehub"
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/lobby"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
	"github.com/gin-gonic/gin"
)

type randomCodeGenerator struct{}

func (r randomCodeGenerator) GenerateCode() string {
	return "AAA"
}

type dummyDispatcher struct{}

func (d dummyDispatcher) Dispatch(m messagebus.Message) {}

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

type PartyService struct {
	hub gamehub.GameHub
}

func (p PartyService) CreateParty() string {
	return p.hub.CreateParty()
}

func (p PartyService) JoinParty(code string, name string) {
	p.hub.HandleMessage(messagebus.JoinParty{Party: messagebus.Party{Code: code}, User: name})
}

func main() {
	hub := gamehub.New(randomCodeGenerator{}, dummyDispatcher{}, randomAllegianceGenerator{})

	router := gin.Default()

	lobby.Register(router, PartyService{hub: hub})

	port := ":44324"
	log.Printf("serving %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("error serving %v\n", err)
	}
}
