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

func newlyConfirmedTeam(service service, code string) gamerules.Game {
	service.handleMessage(joinParty{party: party{code: code}, user: "Alice"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Bob"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Charlie"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Dan"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Edith"})
	service.handleMessage(startGame{party: party{code: code}})
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Alice"})
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Bob"})
	service.handleMessage(leaderConfirmsTeamSelection{party: party{code: code}, leader: "Alice"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _ = game.Start(spiesFirstGenerator{})
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	return game
}

func newlyConductingMission(service service, code string) gamerules.Game {
	service.handleMessage(joinParty{party: party{code: code}, user: "Alice"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Bob"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Charlie"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Dan"})
	service.handleMessage(joinParty{party: party{code: code}, user: "Edith"})
	service.handleMessage(startGame{party: party{code: code}})
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Alice"})
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Bob"})
	service.handleMessage(leaderConfirmsTeamSelection{party: party{code: code}, leader: "Alice"})
	service.handleMessage(approveTeam{party: party{code: code}, player: "Alice"})
	service.handleMessage(approveTeam{party: party{code: code}, player: "Bob"})
	service.handleMessage(approveTeam{party: party{code: code}, player: "Charlie"})
	service.handleMessage(approveTeam{party: party{code: code}, player: "Dan"})
	service.handleMessage(approveTeam{party: party{code: code}, player: "Edith"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _ = game.Start(spiesFirstGenerator{})
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.ApproveTeamBy("Alice")
	game, _ = game.ApproveTeamBy("Bob")
	game, _ = game.ApproveTeamBy("Charlie")
	game, _ = game.ApproveTeamBy("Dan")
	game, _ = game.ApproveTeamBy("Edith")
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

func Test_HandleJoinPartyCommand_IgnoreIfInvalid(t *testing.T) {
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

func Test_HandleStartGameCommand_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(startGame{party: party{code: code}})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(leaderSelectedMember{party: party{code: code}, selectedMember: "Charlie"}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember_IgnoreIfInvalid(t *testing.T) {
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

func Test_HandleLeaderSelectsMember_IgnoreIfWrongLeader(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Bob", memberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsMember(t *testing.T) {
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

func Test_HandleLeaderDeselectsMember_IgnoreIfInvalid(t *testing.T) {
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

func Test_HandleLeaderDeselectsMember_IgnoreIfWrongLeader(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(leaderDeselectsMember{party: party{code: code}, leader: "Bob", memberToDeselect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Dan"})
	service.handleMessage(leaderConfirmsTeamSelection{party: party{code: code}, leader: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(leaderConfirmedSelection{party: party{code: code}}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderSelectsMember("Dan")
	expectedGame, _ = expectedGame.LeaderConfirmsTeamSelection()
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(leaderConfirmsTeamSelection{party: party{code: code}, leader: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection_IgnoreWrongLeader(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Charlie"})
	service.handleMessage(leaderSelectsMember{party: party{code: code}, leader: "Alice", memberToSelect: "Dan"})

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(leaderConfirmsTeamSelection{party: party{code: code}, leader: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderSelectsMember("Dan")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleApproveTeam(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.handleMessage(approveTeam{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(playerVotedOnTeam{party: party{code: code}, player: "Alice", approved: true}))

	expectedGame, _ = expectedGame.ApproveTeamBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(approveTeam{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleRejectTeam(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.handleMessage(rejectTeam{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(playerVotedOnTeam{party: party{code: code}, player: "Alice", approved: false}))

	expectedGame, _ = expectedGame.RejectTeamBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(rejectTeam{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleSucceedMission(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConductingMission(service, code)

	service.handleMessage(succeedMission{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(playerWorkedOnMission{party: party{code: code}, player: "Alice", success: true}))

	expectedGame, _ = expectedGame.SucceedMissionBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(succeedMission{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleFailMission(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConductingMission(service, code)

	service.handleMessage(failMission{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(playerWorkedOnMission{party: party{code: code}, player: "Alice", success: false}))

	expectedGame, _ = expectedGame.FailMissionBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleFailMission_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	messageDispatcher.clearReceivedMessages()
	service.handleMessage(succeedMission{party: party{code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}