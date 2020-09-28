package service

import (
	"testing"

	"github.com/damien-springuel/bomb-canary/server/gamerules"
	. "github.com/onsi/gomega"
)

type testGenerator struct {
	returnCode string
}

func (t testGenerator) generateCode() string {
	return t.returnCode
}

type testMessageDispatcher struct {
	receivedMessage message
}

func (t *testMessageDispatcher) dispatchMessage(m message) {
	t.receivedMessage = m
}

func Test_CreateParty(t *testing.T) {

	s := newService(testGenerator{returnCode: "testCode"}, &testMessageDispatcher{})

	code := s.createParty()

	g := NewWithT(t)
	g.Expect(code).To(Equal("testCode"))
}

func Test_GetGameForPartyCode(t *testing.T) {
	s := newService(testGenerator{returnCode: "testCode"}, &testMessageDispatcher{})

	code := s.createParty()
	game := s.getGameForPartyCode(code)

	expectedGame := gamerules.NewGame()
	g := NewWithT(t)
	g.Expect(game).To(Equal(expectedGame))
}

func Test_HandleJoinPartyCommand(t *testing.T) {

	messageDispatcher := &testMessageDispatcher{}
	s := newService(testGenerator{returnCode: "testCode"}, messageDispatcher)
	code := s.createParty()
	err := s.handleMessage(joinParty{partyCode: "testCode", user: "Alice"})

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(messageDispatcher.receivedMessage).To(Equal(playerJoined{partyCode: "testCode", user: "Alice"}))

	expectedGame := gamerules.NewGame()
	expectedGame, _ = expectedGame.AddPlayer("Alice")
	g.Expect(s.getGameForPartyCode(code)).To(Equal(expectedGame))
}
