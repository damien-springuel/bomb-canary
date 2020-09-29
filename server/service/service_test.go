package service

import (
	"testing"

	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type testGenerator struct {
	returnCode string
}

func (t testGenerator) generateCode() string {
	return t.returnCode
}

type testMessageDispatcher struct {
	lastReceivedMessage messagebus.Message
}

func (t *testMessageDispatcher) dispatchMessage(m messagebus.Message) {
	t.lastReceivedMessage = m
}

type spiesFirstGenerator struct{}

func (s spiesFirstGenerator) Generate(nbPlayers, nbSpies int) []gamerules.Allegiance {
	allegiances := make([]gamerules.Allegiance, nbPlayers)
	for i := range allegiances {
		if i < nbSpies {
			allegiances[i] = gamerules.Spy
		} else {
			allegiances[i] = gamerules.Resistance
		}
	}
	return allegiances
}

func setupService() (*testMessageDispatcher, service, string) {
	messageDispatcher := &testMessageDispatcher{}
	s := newService(testGenerator{returnCode: "testCode"}, messageDispatcher, spiesFirstGenerator{})
	code := s.createParty()
	return messageDispatcher, s, code
}

func Test_CreateParty(t *testing.T) {
	s := newService(testGenerator{returnCode: "testCode"}, nil, nil)
	code := s.createParty()

	g := NewWithT(t)
	g.Expect(code).To(Equal("testCode"))
}

func Test_GetGameForPartyCode(t *testing.T) {
	s := newService(testGenerator{returnCode: "testCode"}, nil, nil)
	code := s.createParty()
	game := s.getGameForPartyCode(code)

	expectedGame := gamerules.NewGame()
	g := NewWithT(t)
	g.Expect(game).To(Equal(expectedGame))
}

func Test_HandleJoinPartyCommand(t *testing.T) {
	messageDispatcher, service, code := setupService()
	service.handleMessage(joinParty{party: party{code: code}, user: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastReceivedMessage).To(Equal(playerJoined{party: party{code: code}, user: "Alice"}))

	expectedGame := gamerules.NewGame()
	expectedGame, _ = expectedGame.AddPlayer("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleStartGameCommand(t *testing.T) {
	messageDispatcher, service, code := setupService()
	service.handleMessage(joinParty{party: party{code: code}, user: "Alice"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Bob"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Charlie"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Dan"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Edith"})
	service.handleMessage(startGame{party: party{code: code}})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastReceivedMessage).To(Equal(gameStarted{party: party{code: code}}))

	expectedGame := gamerules.NewGame()
	expectedGame, _ = expectedGame.AddPlayer("Alice")
	expectedGame, _ = expectedGame.AddPlayer("Bob")
	expectedGame, _ = expectedGame.AddPlayer("Charlie")
	expectedGame, _ = expectedGame.AddPlayer("Dan")
	expectedGame, _ = expectedGame.AddPlayer("Edith")
	expectedGame, _ = expectedGame.Start(spiesFirstGenerator{})
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}
