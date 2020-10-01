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
	g.Expect(dispatcher.receivedMessage).To(Equal(messagebus.StartGame{Party: messagebus.Party{Code: "testCode"}}))
}
