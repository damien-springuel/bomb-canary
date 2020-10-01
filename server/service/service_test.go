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

func (t testGenerator) GenerateCode() string {
	return t.returnCode
}

type testMessageDispatcher struct {
	receivedMessages []messagebus.Message
}

func (t *testMessageDispatcher) Dispatch(m messagebus.Message) {
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

func setupService() (*testMessageDispatcher, Service, string) {
	messageDispatcher := &testMessageDispatcher{}
	s := New(testGenerator{returnCode: "testCode"}, messageDispatcher, spiesFirstGenerator{})
	code := s.CreateParty()
	return messageDispatcher, s, code
}

func newlyStartedGame(service Service, code string) gamerules.Game {
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Alice"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Bob"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Charlie"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Dan"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Edith"})
	service.HandleMessage(StartGame{Party: Party{Code: code}})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _ = game.Start(spiesFirstGenerator{})
	return game
}

func newlyConfirmedTeam(service Service, code string) gamerules.Game {
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Alice"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Bob"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Charlie"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Dan"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Edith"})
	service.HandleMessage(StartGame{Party: Party{Code: code}})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Alice"})

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

func fiveFailedVoteInARow(service Service, code string) gamerules.Game {
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Alice"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Bob"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Charlie"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Dan"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Edith"})
	service.HandleMessage(StartGame{Party: Party{Code: code}})

	// #1
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Edith"})

	// #2
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Edith"})

	// #3
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Charlie", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Charlie", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Edith"})

	// #4
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Dan", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Dan", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Dan"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Edith"})

	// #5 minus Edith's vote
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Edith", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Edith", MemberToSelect: "Edith"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Edith"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Dan"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _ = game.Start(spiesFirstGenerator{})

	// #1
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.RejectTeamBy("Alice")
	game, _ = game.RejectTeamBy("Bob")
	game, _ = game.RejectTeamBy("Charlie")
	game, _ = game.RejectTeamBy("Dan")
	game, _ = game.RejectTeamBy("Edith")

	// #2
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.RejectTeamBy("Alice")
	game, _ = game.RejectTeamBy("Bob")
	game, _ = game.RejectTeamBy("Charlie")
	game, _ = game.RejectTeamBy("Dan")
	game, _ = game.RejectTeamBy("Edith")

	// #3
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.RejectTeamBy("Alice")
	game, _ = game.RejectTeamBy("Bob")
	game, _ = game.RejectTeamBy("Charlie")
	game, _ = game.RejectTeamBy("Dan")
	game, _ = game.RejectTeamBy("Edith")

	// #4
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.RejectTeamBy("Alice")
	game, _ = game.RejectTeamBy("Bob")
	game, _ = game.RejectTeamBy("Charlie")
	game, _ = game.RejectTeamBy("Dan")
	game, _ = game.RejectTeamBy("Edith")

	// #5 minus Edith's vote
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.RejectTeamBy("Alice")
	game, _ = game.RejectTeamBy("Bob")
	game, _ = game.RejectTeamBy("Charlie")
	game, _ = game.RejectTeamBy("Dan")
	return game
}

func newlyConductingMission(service Service, code string) gamerules.Game {
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Alice"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Bob"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Charlie"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Dan"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Edith"})
	service.HandleMessage(StartGame{Party: Party{Code: code}})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})

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

func almostThreeSuccessfulMissions(service Service, code string) gamerules.Game {
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Alice"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Bob"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Charlie"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Dan"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Edith"})
	service.HandleMessage(StartGame{Party: Party{Code: code}})

	// #1
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Bob"})

	// #2
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Charlie"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Charlie"})

	// #3 minus last two succeed
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Charlie", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Charlie", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _ = game.Start(spiesFirstGenerator{})

	// #1
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.ApproveTeamBy("Alice")
	game, _ = game.ApproveTeamBy("Bob")
	game, _ = game.ApproveTeamBy("Charlie")
	game, _ = game.ApproveTeamBy("Dan")
	game, _ = game.ApproveTeamBy("Edith")
	game, _ = game.SucceedMissionBy("Alice")
	game, _ = game.SucceedMissionBy("Bob")

	// #2
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderSelectsMember("Charlie")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.ApproveTeamBy("Alice")
	game, _ = game.ApproveTeamBy("Bob")
	game, _ = game.ApproveTeamBy("Charlie")
	game, _ = game.ApproveTeamBy("Dan")
	game, _ = game.ApproveTeamBy("Edith")
	game, _ = game.SucceedMissionBy("Alice")
	game, _ = game.SucceedMissionBy("Bob")
	game, _ = game.SucceedMissionBy("Charlie")

	// #3
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

func almostThreeFailedMissions(service Service, code string) gamerules.Game {
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Alice"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Bob"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Charlie"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Dan"})
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Edith"})
	service.HandleMessage(StartGame{Party: Party{Code: code}})

	// #1
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Alice"})
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Bob"})

	// #2
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Charlie"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Alice"})
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Bob"})
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Charlie"})

	// #3 minus last two succeed
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Charlie", MemberToSelect: "Alice"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Charlie", MemberToSelect: "Bob"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _ = game.Start(spiesFirstGenerator{})

	// #1
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.ApproveTeamBy("Alice")
	game, _ = game.ApproveTeamBy("Bob")
	game, _ = game.ApproveTeamBy("Charlie")
	game, _ = game.ApproveTeamBy("Dan")
	game, _ = game.ApproveTeamBy("Edith")
	game, _ = game.FailMissionBy("Alice")
	game, _ = game.FailMissionBy("Bob")

	// #2
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderSelectsMember("Charlie")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _ = game.ApproveTeamBy("Alice")
	game, _ = game.ApproveTeamBy("Bob")
	game, _ = game.ApproveTeamBy("Charlie")
	game, _ = game.ApproveTeamBy("Dan")
	game, _ = game.ApproveTeamBy("Edith")
	game, _ = game.FailMissionBy("Alice")
	game, _ = game.FailMissionBy("Bob")
	game, _ = game.FailMissionBy("Charlie")

	// #3
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
	s := New(testGenerator{returnCode: "testCode"}, nil, nil)
	code := s.CreateParty()

	g := NewWithT(t)
	g.Expect(code).To(Equal("testCode"))
}

func Test_GetGameForPartyCode(t *testing.T) {
	s := New(testGenerator{returnCode: "testCode"}, nil, nil)
	code := s.CreateParty()
	game := s.getGameForPartyCode(code)

	expectedGame := gamerules.NewGame()
	g := NewWithT(t)
	g.Expect(game).To(Equal(expectedGame))
}

func Test_HandleJoinPartyCommand(t *testing.T) {
	messageDispatcher, service, code := setupService()
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerJoined{Party: Party{Code: code}, User: "Alice"}))

	expectedGame := gamerules.NewGame()
	expectedGame, _ = expectedGame.AddPlayer("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleJoinPartyCommand_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(JoinParty{Party: Party{Code: code}, User: "Fred"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleShouldIgnoreUnknownMessage(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	type fakeMessage struct{ Party }
	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(fakeMessage{Party: Party{Code: code}})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleStartGameCommand(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(LeaderStartedToSelectMembers{Party: Party{Code: code}, Leader: "Alice"}))
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleStartGameCommand_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(StartGame{Party: Party{Code: code}})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(LeaderSelectedMember{Party: Party{Code: code}, SelectedMember: "Charlie"}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember_IgnoreIfWrongLeader(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsMember(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})
	service.HandleMessage(LeaderDeselectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToDeselect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(LeaderDeselectedMember{Party: Party{Code: code}, DeselectedMember: "Charlie"}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderDeselectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsMember_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(LeaderDeselectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToDeselect: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsMember_IgnoreIfWrongLeader(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(LeaderDeselectsMember{Party: Party{Code: code}, Leader: "Bob", MemberToDeselect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Dan"})
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(LeaderConfirmedSelection{Party: Party{Code: code}}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderSelectsMember("Dan")
	expectedGame, _ = expectedGame.LeaderConfirmsTeamSelection()
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection_IgnoreWrongLeader(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Charlie"})
	service.HandleMessage(LeaderSelectsMember{Party: Party{Code: code}, Leader: "Alice", MemberToSelect: "Dan"})

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(LeaderConfirmsTeamSelection{Party: Party{Code: code}, Leader: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderSelectsMember("Dan")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleApproveTeam(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Alice", Approved: true}))

	expectedGame, _ = expectedGame.ApproveTeamBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_AllPlayerVoted_Approved(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(MissionStarted{Party: Party{Code: code}}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(AllPlayerVotedOnTeam{Party: Party{Code: code}, Approved: true, VoteFailures: 0}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Edith", Approved: true}))

	expectedGame, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _ = expectedGame.ApproveTeamBy("Bob")
	expectedGame, _ = expectedGame.ApproveTeamBy("Charlie")
	expectedGame, _ = expectedGame.ApproveTeamBy("Dan")
	expectedGame, _ = expectedGame.ApproveTeamBy("Edith")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_AllPlayerVoted_Rejected(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Party: Party{Code: code}, Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(AllPlayerVotedOnTeam{Party: Party{Code: code}, Approved: false, VoteFailures: 1}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Edith", Approved: true}))

	expectedGame, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _ = expectedGame.RejectTeamBy("Bob")
	expectedGame, _ = expectedGame.RejectTeamBy("Charlie")
	expectedGame, _ = expectedGame.RejectTeamBy("Dan")
	expectedGame, _ = expectedGame.ApproveTeamBy("Edith")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_RejectedFiveTimeInARow(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := fiveFailedVoteInARow(service, code)

	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Party: Party{Code: code}, Winner: Spy}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(AllPlayerVotedOnTeam{Party: Party{Code: code}, Approved: false, VoteFailures: 5}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Edith", Approved: true}))

	expectedGame, _ = expectedGame.ApproveTeamBy("Edith")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleRejectTeam(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Alice", Approved: false}))

	expectedGame, _ = expectedGame.RejectTeamBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyStartedGame(service, code)

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_AllPlayerVoted_Approved(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(MissionStarted{Party: Party{Code: code}}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(AllPlayerVotedOnTeam{Party: Party{Code: code}, Approved: true, VoteFailures: 0}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Edith", Approved: false}))

	expectedGame, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _ = expectedGame.ApproveTeamBy("Bob")
	expectedGame, _ = expectedGame.ApproveTeamBy("Charlie")
	expectedGame, _ = expectedGame.ApproveTeamBy("Dan")
	expectedGame, _ = expectedGame.RejectTeamBy("Edith")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_AllPlayerVoted_Rejected(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	service.HandleMessage(ApproveTeam{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Bob"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Charlie"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Dan"})
	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Party: Party{Code: code}, Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(AllPlayerVotedOnTeam{Party: Party{Code: code}, Approved: false, VoteFailures: 1}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Edith", Approved: false}))

	expectedGame, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _ = expectedGame.RejectTeamBy("Bob")
	expectedGame, _ = expectedGame.RejectTeamBy("Charlie")
	expectedGame, _ = expectedGame.RejectTeamBy("Dan")
	expectedGame, _ = expectedGame.RejectTeamBy("Edith")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_RejectedFiveTimeInARow(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := fiveFailedVoteInARow(service, code)

	service.HandleMessage(RejectTeam{Party: Party{Code: code}, Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Party: Party{Code: code}, Winner: Spy}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(AllPlayerVotedOnTeam{Party: Party{Code: code}, Approved: false, VoteFailures: 5}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Party: Party{Code: code}, Player: "Edith", Approved: false}))

	expectedGame, _ = expectedGame.RejectTeamBy("Edith")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleSucceedMission(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConductingMission(service, code)

	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerWorkedOnMission{Party: Party{Code: code}, Player: "Alice", Success: true}))

	expectedGame, _ = expectedGame.SucceedMissionBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_MissionCompleted_Successful(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConductingMission(service, code)

	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Party: Party{Code: code}, Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(MissionCompleted{Party: Party{Code: code}, success: true}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Party: Party{Code: code}, Player: "Bob", Success: true}))

	expectedGame, _ = expectedGame.SucceedMissionBy("Alice")
	expectedGame, _ = expectedGame.SucceedMissionBy("Bob")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_MissionCompleted_Failure(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConductingMission(service, code)

	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Alice"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Party: Party{Code: code}, Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(MissionCompleted{Party: Party{Code: code}, success: false}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Party: Party{Code: code}, Player: "Bob", Success: true}))

	expectedGame, _ = expectedGame.FailMissionBy("Alice")
	expectedGame, _ = expectedGame.SucceedMissionBy("Bob")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_MissionCompleted_ThirdSuccess(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := almostThreeSuccessfulMissions(service, code)

	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Alice"})
	service.HandleMessage(SucceedMission{Party: Party{Code: code}, Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Party: Party{Code: code}, Winner: Resistance}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(MissionCompleted{Party: Party{Code: code}, success: true}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Party: Party{Code: code}, Player: "Bob", Success: true}))

	expectedGame, _ = expectedGame.SucceedMissionBy("Alice")
	expectedGame, _ = expectedGame.SucceedMissionBy("Bob")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleFailMission(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConductingMission(service, code)

	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerWorkedOnMission{Party: Party{Code: code}, Player: "Alice", Success: false}))

	expectedGame, _ = expectedGame.FailMissionBy("Alice")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleFailMission_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConfirmedTeam(service, code)

	messageDispatcher.clearReceivedMessages()
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleFailMission_MissionCompleted_Failure(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := newlyConductingMission(service, code)

	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Alice"})
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Party: Party{Code: code}, Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(MissionCompleted{Party: Party{Code: code}, success: false}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Party: Party{Code: code}, Player: "Bob", Success: false}))

	expectedGame, _ = expectedGame.FailMissionBy("Alice")
	expectedGame, _ = expectedGame.FailMissionBy("Bob")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}

func Test_HandleFailMission_MissionCompleted_ThirdFailure(t *testing.T) {
	messageDispatcher, service, code := setupService()
	expectedGame := almostThreeFailedMissions(service, code)

	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Alice"})
	service.HandleMessage(FailMission{Party: Party{Code: code}, player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Party: Party{Code: code}, Winner: Spy}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(MissionCompleted{Party: Party{Code: code}, success: false}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Party: Party{Code: code}, Player: "Bob", Success: false}))

	expectedGame, _ = expectedGame.FailMissionBy("Alice")
	expectedGame, _ = expectedGame.FailMissionBy("Bob")
	g.Expect(service.getGameForPartyCode(code)).To(Equal(expectedGame))
}
