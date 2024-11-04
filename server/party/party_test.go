package party

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

func Test_ServiceJoinParty(t *testing.T) {
	dispatcher := &mockDispatcher{}
	service := NewPartyService(dispatcher)

	g := NewWithT(t)
	service.JoinParty("name")
	g.Expect(dispatcher.receivedMessage).To(Equal(messagebus.JoinParty{Player: "name"}))
}
