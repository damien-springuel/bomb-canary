package playeractions

import (
	"testing"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type mockDispatcher struct {
	receivedMessage messagebus.Message
}

func (m *mockDispatcher) Dispatch(message messagebus.Message) {
	m.receivedMessage = message
}

func Test_ServiceStartGame(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.StartGame("testCode")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(messagebus.StartGame{Command: messagebus.Command{Party: messagebus.Party{Code: "testCode"}}}))
}

func Test_ServiceLeaderSelectsMember(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.LeaderSelectsMember("testCode", "testLeader", "testMember")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.LeaderSelectsMember{
			Command:        messagebus.Command{Party: messagebus.Party{Code: "testCode"}},
			Leader:         "testLeader",
			MemberToSelect: "testMember",
		},
	))
}

func Test_ServiceLeaderDeselectsMember(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.LeaderDeselectsMember("testCode", "testLeader", "testMember")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.LeaderDeselectsMember{
			Command:          messagebus.Command{Party: messagebus.Party{Code: "testCode"}},
			Leader:           "testLeader",
			MemberToDeselect: "testMember",
		},
	))
}

func Test_ServiceLeaderConfirmsTeam(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.LeaderConfirmsTeam("testCode", "testLeader")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.LeaderConfirmsTeamSelection{
			Command: messagebus.Command{Party: messagebus.Party{Code: "testCode"}},
			Leader:  "testLeader",
		},
	))
}

func Test_ServiceApproveTeam(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.ApproveTeam("testCode", "testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.ApproveTeam{
			Command: messagebus.Command{Party: messagebus.Party{Code: "testCode"}},
			Player:  "testPlayer",
		},
	))
}

func Test_ServiceRejectTeam(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.RejectTeam("testCode", "testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.RejectTeam{
			Command: messagebus.Command{Party: messagebus.Party{Code: "testCode"}},
			Player:  "testPlayer",
		},
	))
}

func Test_ServiceSucceedMission(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.SucceedMission("testCode", "testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.SucceedMission{
			Command: messagebus.Command{Party: messagebus.Party{Code: "testCode"}},
			Player:  "testPlayer",
		},
	))
}

func Test_ServiceFailMission(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.FailMission("testCode", "testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.FailMission{
			Command: messagebus.Command{Party: messagebus.Party{Code: "testCode"}},
			Player:  "testPlayer",
		},
	))
}
