package clientstream

import (
	"testing"

	mb "github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type mockEventSender struct {
	receivedCode    string
	receivedMessage []byte

	receivedCodeToPlayer     string
	receivedNameToPlayer     string
	receivedAllNamesToPlayer map[string][]byte
	receivedMessageToPlayer  []byte

	receivedCodeToAllButPlayer     string
	receivedNameToAllButPlayer     string
	receivedMessageToAllButPlayers []byte

	shouldTrackAll      bool
	allReceivedMessages [][]byte
}

func (m *mockEventSender) clearAllReceivedMessages() {
	m.allReceivedMessages = nil
}

func (m *mockEventSender) Send(code string, message []byte) {
	m.receivedCode = code
	m.receivedMessage = message
	if m.shouldTrackAll {
		m.allReceivedMessages = append(m.allReceivedMessages, message)
	}
}

func (m *mockEventSender) SendToPlayer(code, name string, message []byte) {
	m.receivedCodeToPlayer = code
	m.receivedNameToPlayer = name
	m.receivedMessageToPlayer = message
	if m.shouldTrackAll {
		if m.receivedAllNamesToPlayer == nil {
			m.receivedAllNamesToPlayer = make(map[string][]byte)
		}
		m.receivedAllNamesToPlayer[name] = message
		m.allReceivedMessages = append(m.allReceivedMessages, message)
	}
}

func (m *mockEventSender) SendToAllButPlayer(code string, name string, message []byte) {
	m.receivedCodeToAllButPlayer = code
	m.receivedNameToAllButPlayer = name
	m.receivedMessageToAllButPlayers = message
	if m.shouldTrackAll {
		m.allReceivedMessages = append(m.allReceivedMessages, message)
	}
}

func Test_ClientEventBroker_PlayerJoined(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerJoined{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Player: "testName"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{PlayerJoined: &playerJoined{Name: "testName"}}),
		},
	))
}

func Test_ClientEventBroker_PartyCreated(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PartyCreated{Event: mb.Event{Party: mb.Party{Code: "testCode"}}})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{PartyCreated: &partyCreated{Code: "testCode"}}),
		},
	))
}

func Test_ClientEventBroker_PlayerConnected(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerConnected{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Player: "testName"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{PlayerConnected: &playerConnected{Name: "testName"}}),
		},
	))
}

func Test_ClientEventBroker_PlayerDisconnected(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerDisconnected{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Player: "testName"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{PlayerDisconnected: &playerDisconnected{Name: "testName"}}),
		},
	))
}

func Test_ClientEventBroker_GameStarted(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.GameStarted{Event: mb.Event{Party: mb.Party{Code: "testCode"}},
		MissionRequirements: []mb.MissionRequirement{
			{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 2},
			{NbPeopleOnMission: 5, NbFailuresRequiredToFail: 1},
		}})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode: "testCode",
			receivedMessage: toJsonBytes(clientEvent{GameStarted: &gameStarted{
				MissionRequirements: []missionRequirement{
					{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 2},
					{NbPeopleOnMission: 5, NbFailuresRequiredToFail: 1},
				}}}),
		},
	))
}

func Test_ClientEventBroker_SpiesRevealed(t *testing.T) {
	eventSender := &mockEventSender{shouldTrackAll: true}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.AllegianceRevealed{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, AllegianceByPlayer: map[string]mb.Allegiance{
		"p1": mb.Spy,
		"p2": mb.Spy,
		"p3": mb.Resistance,
		"p4": mb.Resistance,
		"p5": mb.Resistance,
	}})

	g := NewWithT(t)
	g.Expect(eventSender.receivedCodeToPlayer).To(Equal("testCode"))
	g.Expect(eventSender.receivedAllNamesToPlayer).To(Equal(
		map[string][]byte{
			"p1": toJsonBytes(clientEvent{SpiesRevealed: &spiesRevealed{Spies: map[string]struct{}{"p1": {}, "p2": {}}}}),
			"p2": toJsonBytes(clientEvent{SpiesRevealed: &spiesRevealed{Spies: map[string]struct{}{"p1": {}, "p2": {}}}}),
			"p3": toJsonBytes(clientEvent{SpiesRevealed: &spiesRevealed{}}),
			"p4": toJsonBytes(clientEvent{SpiesRevealed: &spiesRevealed{}}),
			"p5": toJsonBytes(clientEvent{SpiesRevealed: &spiesRevealed{}}),
		},
	))
}

func Test_ClientEventBroker_LeaderStartedToSelectMembers(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.LeaderStartedToSelectMembers{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Leader: "testLeader"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{LeaderStartedToSelectMembers: &leaderStartedToSelectMembers{Leader: "testLeader"}}),
		},
	))
}

func Test_ClientEventBroker_LeaderSelectedMember(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.LeaderSelectedMember{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, SelectedMember: "testSelectedMember"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{LeaderSelectedMember: &leaderSelectedMember{SelectedMember: "testSelectedMember"}}),
		},
	))
}

func Test_ClientEventBroker_LeaderDeselectedMember(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.LeaderDeselectedMember{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, DeselectedMember: "testDeselectedMember"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{LeaderDeselectedMember: &leaderDeselectedMember{DeselectedMember: "testDeselectedMember"}}),
		},
	))
}

func Test_ClientEventBroker_LeaderConfirmedSelection(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.LeaderConfirmedSelection{Event: mb.Event{Party: mb.Party{Code: "testCode"}}})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{LeaderConfirmedSelection: &leaderConfirmedSelection{}}),
		},
	))
}

func Test_ClientEventBroker_PlayerVotedOnTeam(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerVotedOnTeam{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Player: "testPlayer", Approved: true})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCodeToPlayer:           "testCode",
			receivedNameToPlayer:           "testPlayer",
			receivedMessageToPlayer:        toJsonBytes(clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: "testPlayer", Approved: boolP(true)}}),
			receivedCodeToAllButPlayer:     "testCode",
			receivedNameToAllButPlayer:     "testPlayer",
			receivedMessageToAllButPlayers: toJsonBytes(clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: "testPlayer"}}),
		},
	))
}

func Test_ClientEventBroker_AllPlayerVotedOnTeam(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.AllPlayerVotedOnTeam{
		Event:        mb.Event{Party: mb.Party{Code: "testCode"}},
		Approved:     true,
		VoteFailures: 3,
		PlayerVotes:  map[string]bool{"p1": true, "p2": false},
	})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode: "testCode",
			receivedMessage: toJsonBytes(clientEvent{
				AllPlayerVotedOnTeam: &allPlayerVotedOnTeam{
					Approved:     true,
					VoteFailures: 3,
					PlayerVotes:  map[string]bool{"p1": true, "p2": false},
				},
			}),
		},
	))
}

func Test_ClientEventBroker_MissionStarted(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.MissionStarted{Event: mb.Event{Party: mb.Party{Code: "testCode"}}})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{MissionStarted: &missionStarted{}}),
		},
	))
}

func Test_ClientEventBroker_PlayerWorkedOnMission(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerWorkedOnMission{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Player: "testPlayer", Success: true})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCodeToPlayer:           "testCode",
			receivedNameToPlayer:           "testPlayer",
			receivedMessageToPlayer:        toJsonBytes(clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: "testPlayer", Success: boolP(true)}}),
			receivedCodeToAllButPlayer:     "testCode",
			receivedNameToAllButPlayer:     "testPlayer",
			receivedMessageToAllButPlayers: toJsonBytes(clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: "testPlayer"}}),
		},
	))
}

func Test_ClientEventBroker_MissionCompleted(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(
		mb.MissionCompleted{
			Event:    mb.Event{Party: mb.Party{Code: "testCode"}},
			Success:  true,
			Outcomes: map[bool]int{true: 4, false: 2},
		},
	)

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{MissionCompleted: &missionCompleted{Success: true, NbFails: 2}}),
		},
	))
}

func Test_ClientEventBroker_GameEnded(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.GameEnded{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Winner: mb.Resistance})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{GameEnded: &gameEnded{Winner: string(mb.Resistance)}}),
		},
	))
}
