package clientstream

import (
	"testing"

	mb "github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type mockEventSender struct {
	receivedMessage []byte

	receivedNameToPlayer     string
	receivedAllNamesToPlayer map[string][]byte
	receivedMessageToPlayer  []byte

	receivedNameToAllButPlayer     string
	receivedMessageToAllButPlayers []byte

	shouldTrackAll      bool
	allReceivedMessages [][]byte
}

func (m *mockEventSender) clearAllReceivedMessages() {
	m.allReceivedMessages = nil
}

func (m *mockEventSender) Send(message []byte) {
	m.receivedMessage = message
	if m.shouldTrackAll {
		m.allReceivedMessages = append(m.allReceivedMessages, message)
	}
}

func (m *mockEventSender) SendToPlayer(name string, message []byte) {
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

func (m *mockEventSender) SendToAllButPlayer(name string, message []byte) {
	m.receivedNameToAllButPlayer = name
	m.receivedMessageToAllButPlayers = message
	if m.shouldTrackAll {
		m.allReceivedMessages = append(m.allReceivedMessages, message)
	}
}

func Test_ClientEventBroker_PlayerJoined(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerJoined{Player: "testName"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{PlayerJoined: &playerJoined{Name: "testName"}}),
		},
	))
}

func Test_ClientEventBroker_PlayerConnected(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerConnected{Player: "testName"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{PlayerConnected: &playerConnected{Name: "testName"}}),
		},
	))
}

func Test_ClientEventBroker_PlayerDisconnected(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerDisconnected{Player: "testName"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{PlayerDisconnected: &playerDisconnected{Name: "testName"}}),
		},
	))
}

func Test_ClientEventBroker_GameStarted(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.GameStarted{
		MissionRequirements: []mb.MissionRequirement{
			{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 2},
			{NbPeopleOnMission: 5, NbFailuresRequiredToFail: 1},
		}})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
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
	eventBroker.Consume(mb.AllegianceRevealed{AllegianceByPlayer: map[string]mb.Allegiance{
		"p1": mb.Spy,
		"p2": mb.Spy,
		"p3": mb.Resistance,
		"p4": mb.Resistance,
		"p5": mb.Resistance,
	}})

	g := NewWithT(t)
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
	eventBroker.Consume(mb.LeaderStartedToSelectMembers{Leader: "testLeader"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{LeaderStartedToSelectMembers: &leaderStartedToSelectMembers{Leader: "testLeader"}}),
		},
	))
}

func Test_ClientEventBroker_LeaderSelectedMember(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.LeaderSelectedMember{SelectedMember: "testSelectedMember"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{LeaderSelectedMember: &leaderSelectedMember{SelectedMember: "testSelectedMember"}}),
		},
	))
}

func Test_ClientEventBroker_LeaderDeselectedMember(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.LeaderDeselectedMember{DeselectedMember: "testDeselectedMember"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{LeaderDeselectedMember: &leaderDeselectedMember{DeselectedMember: "testDeselectedMember"}}),
		},
	))
}

func Test_ClientEventBroker_LeaderConfirmedSelection(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.LeaderConfirmedSelection{})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{LeaderConfirmedSelection: &leaderConfirmedSelection{}}),
		},
	))
}

func Test_ClientEventBroker_PlayerVotedOnTeam(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerVotedOnTeam{Player: "testPlayer", Approved: true})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedNameToPlayer:           "testPlayer",
			receivedMessageToPlayer:        toJsonBytes(clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: "testPlayer", Approved: boolP(true)}}),
			receivedNameToAllButPlayer:     "testPlayer",
			receivedMessageToAllButPlayers: toJsonBytes(clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: "testPlayer"}}),
		},
	))
}

func Test_ClientEventBroker_AllPlayerVotedOnTeam(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.AllPlayerVotedOnTeam{
		Approved:     true,
		VoteFailures: 3,
		PlayerVotes:  map[string]bool{"p1": true, "p2": false},
	})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
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
	eventBroker.Consume(mb.MissionStarted{})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{MissionStarted: &missionStarted{}}),
		},
	))
}

func Test_ClientEventBroker_PlayerWorkedOnMission(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerWorkedOnMission{Player: "testPlayer", Success: true})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedNameToPlayer:           "testPlayer",
			receivedMessageToPlayer:        toJsonBytes(clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: "testPlayer", Success: boolP(true)}}),
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
			Success:  true,
			Outcomes: map[bool]int{true: 4, false: 2},
		},
	)

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{MissionCompleted: &missionCompleted{Success: true, NbFails: 2}}),
		},
	))
}

func Test_ClientEventBroker_GameEnded(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.GameEnded{Winner: mb.Resistance, Spies: []string{"p1", "p2"}})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedMessage: toJsonBytes(clientEvent{GameEnded: &gameEnded{Winner: string(mb.Resistance), Spies: []string{"p1", "p2"}}}),
		},
	))
}
