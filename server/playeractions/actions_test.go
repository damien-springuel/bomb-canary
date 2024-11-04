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

	s.StartGame()

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(messagebus.StartGame{}))
}

func Test_ServiceLeaderSelectsMember(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.LeaderSelectsMember("testLeader", "testMember")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.LeaderSelectsMember{
			Leader:         "testLeader",
			MemberToSelect: "testMember",
		},
	))
}

func Test_ServiceLeaderDeselectsMember(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.LeaderDeselectsMember("testLeader", "testMember")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.LeaderDeselectsMember{
			Leader:           "testLeader",
			MemberToDeselect: "testMember",
		},
	))
}

func Test_ServiceLeaderConfirmsTeam(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.LeaderConfirmsTeam("testLeader")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.LeaderConfirmsTeamSelection{
			Leader: "testLeader",
		},
	))
}

func Test_ServiceApproveTeam(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.ApproveTeam("testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.ApproveTeam{
			Player: "testPlayer",
		},
	))
}

func Test_ServiceRejectTeam(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.RejectTeam("testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.RejectTeam{
			Player: "testPlayer",
		},
	))
}

func Test_ServiceSucceedMission(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.SucceedMission("testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.SucceedMission{
			Player: "testPlayer",
		},
	))
}

func Test_ServiceFailMission(t *testing.T) {
	dispatcher := &mockDispatcher{}
	s := NewActionService(dispatcher)

	s.FailMission("testPlayer")

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessage).To(Equal(
		messagebus.FailMission{
			Player: "testPlayer",
		},
	))
}
