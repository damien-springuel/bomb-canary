package gamehub

import (
	"testing"

	"github.com/damien-springuel/bomb-canary/server/gamerules"
	. "github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type testMessageDispatcher struct {
	receivedMessages []Message
}

func (t *testMessageDispatcher) Dispatch(m Message) {
	t.receivedMessages = append(t.receivedMessages, m)
}

func (t *testMessageDispatcher) messageFromEnd(index int) Message {
	return t.receivedMessages[len(t.receivedMessages)-1-index]
}

func (t *testMessageDispatcher) lastMessage() Message {
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

func setupHub() (*testMessageDispatcher, *gameHub) {
	messageDispatcher := &testMessageDispatcher{}
	hub := New(messageDispatcher, spiesFirstGenerator{})
	return messageDispatcher, hub
}

func newlyStartedGame(hub *gameHub) gamerules.Game {
	hub.Consume(JoinParty{Player: "Alice"})
	hub.Consume(JoinParty{Player: "Bob"})
	hub.Consume(JoinParty{Player: "Charlie"})
	hub.Consume(JoinParty{Player: "Dan"})
	hub.Consume(JoinParty{Player: "Edith"})
	hub.Consume(StartGame{})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _, _, _ = game.Start(spiesFirstGenerator{})
	return game
}

func newlyConfirmedTeam(hub *gameHub) gamerules.Game {
	hub.Consume(JoinParty{Player: "Alice"})
	hub.Consume(JoinParty{Player: "Bob"})
	hub.Consume(JoinParty{Player: "Charlie"})
	hub.Consume(JoinParty{Player: "Dan"})
	hub.Consume(JoinParty{Player: "Edith"})
	hub.Consume(StartGame{})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _, _, _ = game.Start(spiesFirstGenerator{})
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	return game
}

func fiveFailedVoteInARow(hub *gameHub) gamerules.Game {
	hub.Consume(JoinParty{Player: "Alice"})
	hub.Consume(JoinParty{Player: "Bob"})
	hub.Consume(JoinParty{Player: "Charlie"})
	hub.Consume(JoinParty{Player: "Dan"})
	hub.Consume(JoinParty{Player: "Edith"})
	hub.Consume(StartGame{})

	// #1
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})
	hub.Consume(RejectTeam{Player: "Alice"})
	hub.Consume(RejectTeam{Player: "Bob"})
	hub.Consume(RejectTeam{Player: "Charlie"})
	hub.Consume(RejectTeam{Player: "Dan"})
	hub.Consume(RejectTeam{Player: "Edith"})

	// #2
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Bob"})
	hub.Consume(RejectTeam{Player: "Alice"})
	hub.Consume(RejectTeam{Player: "Bob"})
	hub.Consume(RejectTeam{Player: "Charlie"})
	hub.Consume(RejectTeam{Player: "Dan"})
	hub.Consume(RejectTeam{Player: "Edith"})

	// #3
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Charlie"})
	hub.Consume(RejectTeam{Player: "Alice"})
	hub.Consume(RejectTeam{Player: "Bob"})
	hub.Consume(RejectTeam{Player: "Charlie"})
	hub.Consume(RejectTeam{Player: "Dan"})
	hub.Consume(RejectTeam{Player: "Edith"})

	// #4
	hub.Consume(LeaderSelectsMember{Leader: "Dan", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Dan", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Dan"})
	hub.Consume(RejectTeam{Player: "Alice"})
	hub.Consume(RejectTeam{Player: "Bob"})
	hub.Consume(RejectTeam{Player: "Charlie"})
	hub.Consume(RejectTeam{Player: "Dan"})
	hub.Consume(RejectTeam{Player: "Edith"})

	// #5 minus Edith's vote
	hub.Consume(LeaderSelectsMember{Leader: "Edith", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Edith", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Edith"})
	hub.Consume(RejectTeam{Player: "Alice"})
	hub.Consume(RejectTeam{Player: "Bob"})
	hub.Consume(RejectTeam{Player: "Charlie"})
	hub.Consume(RejectTeam{Player: "Dan"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _, _, _ = game.Start(spiesFirstGenerator{})

	// #1
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.RejectTeamBy("Alice")
	game, _, _ = game.RejectTeamBy("Bob")
	game, _, _ = game.RejectTeamBy("Charlie")
	game, _, _ = game.RejectTeamBy("Dan")
	game, _, _ = game.RejectTeamBy("Edith")

	// #2
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.RejectTeamBy("Alice")
	game, _, _ = game.RejectTeamBy("Bob")
	game, _, _ = game.RejectTeamBy("Charlie")
	game, _, _ = game.RejectTeamBy("Dan")
	game, _, _ = game.RejectTeamBy("Edith")

	// #3
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.RejectTeamBy("Alice")
	game, _, _ = game.RejectTeamBy("Bob")
	game, _, _ = game.RejectTeamBy("Charlie")
	game, _, _ = game.RejectTeamBy("Dan")
	game, _, _ = game.RejectTeamBy("Edith")

	// #4
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.RejectTeamBy("Alice")
	game, _, _ = game.RejectTeamBy("Bob")
	game, _, _ = game.RejectTeamBy("Charlie")
	game, _, _ = game.RejectTeamBy("Dan")
	game, _, _ = game.RejectTeamBy("Edith")

	// #5 minus Edith's vote
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.RejectTeamBy("Alice")
	game, _, _ = game.RejectTeamBy("Bob")
	game, _, _ = game.RejectTeamBy("Charlie")
	game, _, _ = game.RejectTeamBy("Dan")
	return game
}

func newlyConductingMission(hub *gameHub) gamerules.Game {
	hub.Consume(JoinParty{Player: "Alice"})
	hub.Consume(JoinParty{Player: "Bob"})
	hub.Consume(JoinParty{Player: "Charlie"})
	hub.Consume(JoinParty{Player: "Dan"})
	hub.Consume(JoinParty{Player: "Edith"})
	hub.Consume(StartGame{})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _, _, _ = game.Start(spiesFirstGenerator{})
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	return game
}

func almostThreeSuccessfulMissions(hub *gameHub) gamerules.Game {
	hub.Consume(JoinParty{Player: "Alice"})
	hub.Consume(JoinParty{Player: "Bob"})
	hub.Consume(JoinParty{Player: "Charlie"})
	hub.Consume(JoinParty{Player: "Dan"})
	hub.Consume(JoinParty{Player: "Edith"})
	hub.Consume(StartGame{})

	// #1
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(SucceedMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})

	// #2
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Bob"})
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Charlie"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Bob"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(SucceedMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})
	hub.Consume(SucceedMission{Player: "Charlie"})

	// #3 minus last two succeed
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _, _, _ = game.Start(spiesFirstGenerator{})

	// #1
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.SucceedMissionBy("Alice")
	game, _, _ = game.SucceedMissionBy("Bob")

	// #2
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderSelectsMember("Charlie")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.SucceedMissionBy("Alice")
	game, _, _ = game.SucceedMissionBy("Bob")
	game, _, _ = game.SucceedMissionBy("Charlie")

	// #3
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	return game
}

func almostThreeFailedMissions(hub *gameHub) gamerules.Game {
	hub.Consume(JoinParty{Player: "Alice"})
	hub.Consume(JoinParty{Player: "Bob"})
	hub.Consume(JoinParty{Player: "Charlie"})
	hub.Consume(JoinParty{Player: "Dan"})
	hub.Consume(JoinParty{Player: "Edith"})
	hub.Consume(StartGame{})

	// #1
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(FailMission{Player: "Alice"})
	hub.Consume(FailMission{Player: "Bob"})

	// #2
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Bob"})
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Charlie"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Bob"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(FailMission{Player: "Alice"})
	hub.Consume(FailMission{Player: "Bob"})
	hub.Consume(FailMission{Player: "Charlie"})

	// #3 minus last two fails
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _, _, _ = game.Start(spiesFirstGenerator{})

	// #1
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.FailMissionBy("Alice")
	game, _, _ = game.FailMissionBy("Bob")

	// #2
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderSelectsMember("Charlie")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.FailMissionBy("Alice")
	game, _, _ = game.FailMissionBy("Bob")
	game, _, _ = game.FailMissionBy("Charlie")

	// #3
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	return game
}

func twoSuccessTwoFailuresLastMission(hub *gameHub) gamerules.Game {
	hub.Consume(JoinParty{Player: "Alice"})
	hub.Consume(JoinParty{Player: "Bob"})
	hub.Consume(JoinParty{Player: "Charlie"})
	hub.Consume(JoinParty{Player: "Dan"})
	hub.Consume(JoinParty{Player: "Edith"})
	hub.Consume(StartGame{})

	// #1
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(SucceedMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})

	// #2
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Bob"})
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Charlie"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Bob"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(FailMission{Player: "Alice"})
	hub.Consume(FailMission{Player: "Bob"})
	hub.Consume(FailMission{Player: "Charlie"})

	// #3
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Charlie", MemberToSelect: "Bob"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(SucceedMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})

	// #4
	hub.Consume(LeaderSelectsMember{Leader: "Dan", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Dan", MemberToSelect: "Bob"})
	hub.Consume(LeaderSelectsMember{Leader: "Dan", MemberToSelect: "Charlie"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Dan"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})
	hub.Consume(FailMission{Player: "Alice"})
	hub.Consume(FailMission{Player: "Bob"})
	hub.Consume(FailMission{Player: "Charlie"})

	// #5
	hub.Consume(LeaderSelectsMember{Leader: "Edith", MemberToSelect: "Alice"})
	hub.Consume(LeaderSelectsMember{Leader: "Edith", MemberToSelect: "Bob"})
	hub.Consume(LeaderSelectsMember{Leader: "Edith", MemberToSelect: "Charlie"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Edith"})
	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})

	game := gamerules.NewGame()
	game, _ = game.AddPlayer("Alice")
	game, _ = game.AddPlayer("Bob")
	game, _ = game.AddPlayer("Charlie")
	game, _ = game.AddPlayer("Dan")
	game, _ = game.AddPlayer("Edith")
	game, _, _, _ = game.Start(spiesFirstGenerator{})

	// #1
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.SucceedMissionBy("Alice")
	game, _, _ = game.SucceedMissionBy("Bob")

	// #2
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderSelectsMember("Charlie")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.FailMissionBy("Alice")
	game, _, _ = game.FailMissionBy("Bob")
	game, _, _ = game.FailMissionBy("Charlie")

	// #3
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.SucceedMissionBy("Alice")
	game, _, _ = game.SucceedMissionBy("Bob")

	// #4
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderSelectsMember("Charlie")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	game, _, _ = game.FailMissionBy("Alice")
	game, _, _ = game.FailMissionBy("Bob")
	game, _, _ = game.FailMissionBy("Charlie")

	// #5
	game, _ = game.LeaderSelectsMember("Alice")
	game, _ = game.LeaderSelectsMember("Bob")
	game, _ = game.LeaderSelectsMember("Charlie")
	game, _ = game.LeaderConfirmsTeamSelection()
	game, _, _ = game.ApproveTeamBy("Alice")
	game, _, _ = game.ApproveTeamBy("Bob")
	game, _, _ = game.ApproveTeamBy("Charlie")
	game, _, _ = game.ApproveTeamBy("Dan")
	game, _, _ = game.ApproveTeamBy("Edith")
	return game
}

func Test_HandleJoinPartyCommand(t *testing.T) {
	messageDispatcher, hub := setupHub()
	hub.Consume(JoinParty{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerJoined{Player: "Alice"}))

	expectedGame := gamerules.NewGame()
	expectedGame, _ = expectedGame.AddPlayer("Alice")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleJoinPartyCommand_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)

	messageDispatcher.clearReceivedMessages()
	hub.Consume(JoinParty{Player: "Fred"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleShouldIgnoreUnknownMessage(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)

	type fakeMessage struct{ Command }
	messageDispatcher.clearReceivedMessages()
	hub.Consume(fakeMessage{})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleStartGameCommand(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Leader: "Alice"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(AllegianceRevealed{AllegianceByPlayer: map[string]Allegiance{
		"Alice":   Spy,
		"Bob":     Spy,
		"Charlie": Resistance,
		"Dan":     Resistance,
		"Edith":   Resistance,
	}}))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(GameStarted{MissionRequirements: []MissionRequirement{
		{NbPeopleOnMission: 2, NbFailuresRequiredToFail: 1},
		{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 1},
		{NbPeopleOnMission: 2, NbFailuresRequiredToFail: 1},
		{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 1},
		{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 1},
	}}))
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleStartGameCommand_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)

	messageDispatcher.clearReceivedMessages()
	hub.Consume(StartGame{})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(LeaderSelectedMember{SelectedMember: "Charlie"}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderSelectsMember_IgnoreIfWrongLeader(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)

	messageDispatcher.clearReceivedMessages()
	hub.Consume(LeaderSelectsMember{Leader: "Bob", MemberToSelect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsMember(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})
	hub.Consume(LeaderDeselectsMember{Leader: "Alice", MemberToDeselect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(LeaderDeselectedMember{DeselectedMember: "Charlie"}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderDeselectsMember("Charlie")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsMember_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	hub.Consume(LeaderDeselectsMember{Leader: "Alice", MemberToDeselect: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderDeselectsMember_IgnoreIfWrongLeader(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	hub.Consume(LeaderDeselectsMember{Leader: "Bob", MemberToDeselect: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Dan"})
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(LeaderConfirmedSelection{}))

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderSelectsMember("Dan")
	expectedGame, _ = expectedGame.LeaderConfirmsTeamSelection()
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})

	messageDispatcher.clearReceivedMessages()
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleLeaderConfirmsSelection_IgnoreWrongLeader(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Charlie"})
	hub.Consume(LeaderSelectsMember{Leader: "Alice", MemberToSelect: "Dan"})

	messageDispatcher.clearReceivedMessages()
	hub.Consume(LeaderConfirmsTeamSelection{Leader: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeEmpty())

	expectedGame, _ = expectedGame.LeaderSelectsMember("Charlie")
	expectedGame, _ = expectedGame.LeaderSelectsMember("Dan")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleApproveTeam(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	hub.Consume(ApproveTeam{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerVotedOnTeam{Player: "Alice", Approved: true}))

	expectedGame, _, _ = expectedGame.ApproveTeamBy("Alice")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)

	messageDispatcher.clearReceivedMessages()
	hub.Consume(ApproveTeam{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_AllPlayerVoted_Approved(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(MissionStarted{}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		AllPlayerVotedOnTeam{
			Approved:     true,
			VoteFailures: 0,
			PlayerVotes:  map[string]bool{"Alice": true, "Bob": true, "Charlie": true, "Dan": true, "Edith": true},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Player: "Edith", Approved: true}))

	expectedGame, _, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Bob")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Charlie")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Dan")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Edith")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_AllPlayerVoted_Rejected(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(RejectTeam{Player: "Bob"})
	hub.Consume(RejectTeam{Player: "Charlie"})
	hub.Consume(RejectTeam{Player: "Dan"})
	hub.Consume(ApproveTeam{Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		AllPlayerVotedOnTeam{
			Approved:     false,
			VoteFailures: 1,
			PlayerVotes:  map[string]bool{"Alice": true, "Bob": false, "Charlie": false, "Dan": false, "Edith": true},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Player: "Edith", Approved: true}))

	expectedGame, _, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Bob")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Charlie")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Dan")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Edith")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleApproveTeam_RejectedFiveTimeInARow(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := fiveFailedVoteInARow(hub)

	hub.Consume(ApproveTeam{Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Winner: Spy, Spies: []string{"Alice", "Bob"}}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		AllPlayerVotedOnTeam{
			Approved:     false,
			VoteFailures: 5,
			PlayerVotes:  map[string]bool{"Alice": false, "Bob": false, "Charlie": false, "Dan": false, "Edith": true},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Player: "Edith", Approved: true}))

	expectedGame, _, _ = expectedGame.ApproveTeamBy("Edith")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleRejectTeam(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	hub.Consume(RejectTeam{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerVotedOnTeam{Player: "Alice", Approved: false}))

	expectedGame, _, _ = expectedGame.RejectTeamBy("Alice")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyStartedGame(hub)

	messageDispatcher.clearReceivedMessages()
	hub.Consume(RejectTeam{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_AllPlayerVoted_Approved(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(ApproveTeam{Player: "Bob"})
	hub.Consume(ApproveTeam{Player: "Charlie"})
	hub.Consume(ApproveTeam{Player: "Dan"})
	hub.Consume(RejectTeam{Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(MissionStarted{}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		AllPlayerVotedOnTeam{
			Approved:     true,
			VoteFailures: 0,
			PlayerVotes:  map[string]bool{"Alice": true, "Bob": true, "Charlie": true, "Dan": true, "Edith": false},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Player: "Edith", Approved: false}))

	expectedGame, _, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Bob")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Charlie")
	expectedGame, _, _ = expectedGame.ApproveTeamBy("Dan")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Edith")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_AllPlayerVoted_Rejected(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	hub.Consume(ApproveTeam{Player: "Alice"})
	hub.Consume(RejectTeam{Player: "Bob"})
	hub.Consume(RejectTeam{Player: "Charlie"})
	hub.Consume(RejectTeam{Player: "Dan"})
	hub.Consume(RejectTeam{Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		AllPlayerVotedOnTeam{
			Approved:     false,
			VoteFailures: 1,
			PlayerVotes:  map[string]bool{"Alice": true, "Bob": false, "Charlie": false, "Dan": false, "Edith": false},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Player: "Edith", Approved: false}))

	expectedGame, _, _ = expectedGame.ApproveTeamBy("Alice")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Bob")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Charlie")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Dan")
	expectedGame, _, _ = expectedGame.RejectTeamBy("Edith")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleRejectTeam_RejectedFiveTimeInARow(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := fiveFailedVoteInARow(hub)

	hub.Consume(RejectTeam{Player: "Edith"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Winner: Spy, Spies: []string{"Alice", "Bob"}}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		AllPlayerVotedOnTeam{
			Approved:     false,
			VoteFailures: 5,
			PlayerVotes:  map[string]bool{"Alice": false, "Bob": false, "Charlie": false, "Dan": false, "Edith": false},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerVotedOnTeam{Player: "Edith", Approved: false}))

	expectedGame, _, _ = expectedGame.RejectTeamBy("Edith")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleSucceedMission(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConductingMission(hub)

	hub.Consume(SucceedMission{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerWorkedOnMission{Player: "Alice", Success: true}))

	expectedGame, _, _ = expectedGame.SucceedMissionBy("Alice")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	messageDispatcher.clearReceivedMessages()
	hub.Consume(SucceedMission{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_MissionCompleted_Successful(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConductingMission(hub)

	hub.Consume(SucceedMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		MissionCompleted{
			Success:  true,
			Outcomes: map[bool]int{true: 2},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Player: "Bob", Success: true}))

	expectedGame, _, _ = expectedGame.SucceedMissionBy("Alice")
	expectedGame, _, _ = expectedGame.SucceedMissionBy("Bob")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_MissionCompleted_Failure(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConductingMission(hub)

	hub.Consume(FailMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		MissionCompleted{
			Success:  false,
			Outcomes: map[bool]int{false: 1, true: 1},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Player: "Bob", Success: true}))

	expectedGame, _, _ = expectedGame.FailMissionBy("Alice")
	expectedGame, _, _ = expectedGame.SucceedMissionBy("Bob")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_MissionCompleted_ThirdSuccess(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := almostThreeSuccessfulMissions(hub)

	hub.Consume(SucceedMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Winner: Resistance, Spies: []string{"Alice", "Bob"}}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		MissionCompleted{
			Success:  true,
			Outcomes: map[bool]int{true: 2},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Player: "Bob", Success: true}))

	expectedGame, _, _ = expectedGame.SucceedMissionBy("Alice")
	expectedGame, _, _ = expectedGame.SucceedMissionBy("Bob")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleSucceedMission_MissionCompleted_LastMission(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := twoSuccessTwoFailuresLastMission(hub)

	hub.Consume(SucceedMission{Player: "Alice"})
	hub.Consume(SucceedMission{Player: "Bob"})
	hub.Consume(SucceedMission{Player: "Charlie"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Winner: Resistance, Spies: []string{"Alice", "Bob"}}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		MissionCompleted{
			Success:  true,
			Outcomes: map[bool]int{true: 3},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Player: "Charlie", Success: true}))

	expectedGame, _, _ = expectedGame.SucceedMissionBy("Alice")
	expectedGame, _, _ = expectedGame.SucceedMissionBy("Bob")
	expectedGame, _, _ = expectedGame.SucceedMissionBy("Charlie")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleFailMission(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConductingMission(hub)

	hub.Consume(FailMission{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.lastMessage()).To(Equal(PlayerWorkedOnMission{Player: "Alice", Success: false}))

	expectedGame, _, _ = expectedGame.FailMissionBy("Alice")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleFailMission_IgnoreIfInvalid(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConfirmedTeam(hub)

	messageDispatcher.clearReceivedMessages()
	hub.Consume(FailMission{Player: "Alice"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.receivedMessages).To(BeNil())
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleFailMission_MissionCompleted_Failure(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := newlyConductingMission(hub)

	hub.Consume(FailMission{Player: "Alice"})
	hub.Consume(FailMission{Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(LeaderStartedToSelectMembers{Leader: "Bob"}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		MissionCompleted{
			Success:  false,
			Outcomes: map[bool]int{false: 2},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Player: "Bob", Success: false}))

	expectedGame, _, _ = expectedGame.FailMissionBy("Alice")
	expectedGame, _, _ = expectedGame.FailMissionBy("Bob")
	g.Expect(hub.game).To(Equal(expectedGame))
}

func Test_HandleFailMission_MissionCompleted_ThirdFailure(t *testing.T) {
	messageDispatcher, hub := setupHub()
	expectedGame := almostThreeFailedMissions(hub)

	hub.Consume(FailMission{Player: "Alice"})
	hub.Consume(FailMission{Player: "Bob"})

	g := NewWithT(t)
	g.Expect(messageDispatcher.messageFromEnd(0)).To(Equal(GameEnded{Winner: Spy, Spies: []string{"Alice", "Bob"}}))
	g.Expect(messageDispatcher.messageFromEnd(1)).To(Equal(
		MissionCompleted{
			Success:  false,
			Outcomes: map[bool]int{false: 2},
		},
	))
	g.Expect(messageDispatcher.messageFromEnd(2)).To(Equal(PlayerWorkedOnMission{Player: "Bob", Success: false}))

	expectedGame, _, _ = expectedGame.FailMissionBy("Alice")
	expectedGame, _, _ = expectedGame.FailMissionBy("Bob")
	g.Expect(hub.game).To(Equal(expectedGame))
}
