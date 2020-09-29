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
	receivedMessages []messagebus.Message
}

func (t *testMessageDispatcher) dispatchMessage(m messagebus.Message) {
	t.receivedMessages = append(t.receivedMessages, m)
}

func (t *testMessageDispatcher) messageFromEnd(index int) messagebus.Message {
	return t.receivedMessages[len(t.receivedMessages)-1-index]
}

func (t *testMessageDispatcher) lastMessage() messagebus.Message {
	return t.messageFromEnd(0)
}

func (t *testMessageDispatcher) clearReceivedMessages() {
	t.receivedMessages = nil
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

func newlyStartedGame(service service, code string) gamerules.Game {
	service.handleMessage(joinParty{party: party{code: code}, user: "Alice"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Bob"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Charlie"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Dan"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Edith"})
	service.handleMessage(startGame{party: party{code: code}})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _ = game.Start(spiesFirstGenerator{})
	return game
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
	g.Expect(messageDispatcher.lastMessage()).To(Equal(playerJoined{party: party{code: code}, user: "Alice"}))

	expectedGame := gamerules.NewGame()
	expectedGame, _ = expectedGame.AddPlayer("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleJoinPartyCommand_ShouldIgnoreIfGameAlreadyStarted(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(joinParty{party: party{code: code}, user: "Fred"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleShouldIgnoreUnknownMessage(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	type fakeMessage struct{ party }
	messageDispatcher.clearReceivedMessages()
	service.handleMessage(fakeMessage{party: party{code: code}})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleStartGameCommand(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(leaderStartedToSelectMembers{party: party{code: code}, leader: "Alice"}))
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleStartGameCommand_ShouldIgnoreIfGameAlreadyStarted(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(startGame{party: party{code: code}})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsAMember(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(leaderSelectedMember{party: party{code: code}, selectedMember: "Charlie"}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsAMember_ShouldIgnoreIfMemberAlreadySelected(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsAMember(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})
	service.handleMessage(leaderDeselectsMember{party: party{code: code}, leader: "Alice", memberToDeselect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(leaderDeselectedMember{party: party{code: code}, deselectedMember: "Charlie"}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderDeselectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsAMember_ShouldIgnoreIfMemberNotSelected(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(leaderDeselectsMember{party: party{code: code}, leader: "Alice", memberToDeselect: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}
