package clientstream

import (
	"testing"

	mb "github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type mockEventSender struct {
	receivedCode    string
	receivedMessage []byte

	receiveCodeToPlayer     string
	receiveNameToPlayer     string
	receivedMessageToPlayer []byte

	receiveCodeToAllButPlayer     string
	receiveNameToAllButPlayer     string
	receivedMessageToAllButPlayer []byte

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
	m.receiveCodeToPlayer = code
	m.receiveNameToPlayer = name
	m.receivedMessageToPlayer = message
	if m.shouldTrackAll {
		m.allReceivedMessages = append(m.allReceivedMessages, message)
	}
}

func (m *mockEventSender) SendToAllButPlayer(code, name string, message []byte) {
	m.receiveCodeToAllButPlayer = code
	m.receiveNameToAllButPlayer = name
	m.receivedMessageToAllButPlayer = message
	if m.shouldTrackAll {
		m.allReceivedMessages = append(m.allReceivedMessages, message)
	}
}

func Test_ClientEventBroker_PlayerJoined(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerJoined{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Player: "testUser"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{PlayerJoined: &playerJoined{Name: "testUser"}}),
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
			receiveCodeToPlayer:           "testCode",
			receiveNameToPlayer:           "testPlayer",
			receivedMessageToPlayer:       toJsonBytes(clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: "testPlayer", Approved: boolP(true)}}),
			receiveCodeToAllButPlayer:     "testCode",
			receiveNameToAllButPlayer:     "testPlayer",
			receivedMessageToAllButPlayer: toJsonBytes(clientEvent{PlayerVotedOnTeam: &playerVotedOnTeam{Player: "testPlayer"}}),
		},
	))
}

func Test_ClientEventBroker_AllPlayerVotedOnTeam(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.AllPlayerVotedOnTeam{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Approved: true, VoteFailures: 3})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{AllPlayerVotedOnTeam: &allPlayerVotedOnTeam{Approved: true, VoteFailures: 3}}),
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
			receiveCodeToPlayer:           "testCode",
			receiveNameToPlayer:           "testPlayer",
			receivedMessageToPlayer:       toJsonBytes(clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: "testPlayer", Success: boolP(true)}}),
			receiveCodeToAllButPlayer:     "testCode",
			receiveNameToAllButPlayer:     "testPlayer",
			receivedMessageToAllButPlayer: toJsonBytes(clientEvent{PlayerWorkedOnMission: &playerWorkedOnMission{Player: "testPlayer"}}),
		},
	))
}

func Test_ClientEventBroker_MissionCompleted(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.MissionCompleted{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Success: true})

	g := NewWithT(t)
	g.Expect(eventSender.receivedCode).To(Equal("testCode"))
	g.Expect(eventSender.receivedCode).To(Equal("testCode"))
	g.Expect(*eventSender).To(Equal(
		mockEventSender{
			receivedCode:    "testCode",
			receivedMessage: toJsonBytes(clientEvent{MissionCompleted: &missionCompleted{Success: true}}),
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
