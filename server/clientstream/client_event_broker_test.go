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
}

func (m *mockEventSender) Send(code string, message []byte) {
	m.receivedCode = code
	m.receivedMessage = message
}

func (m *mockEventSender) SendToPlayer(code, name string, message []byte) {
	m.receiveCodeToPlayer = code
	m.receiveNameToPlayer = name
	m.receivedMessageToPlayer = message
}

func (m *mockEventSender) SendToAllButPlayer(code, name string, message []byte) {
	m.receiveCodeToAllButPlayer = code
	m.receiveNameToAllButPlayer = name
	m.receivedMessageToAllButPlayer = message
}

func Test_PlayerJoined(t *testing.T) {
	eventSender := &mockEventSender{}
	eventBroker := NewClientEventBroker(eventSender)
	eventBroker.Consume(mb.PlayerJoined{Event: mb.Event{Party: mb.Party{Code: "testCode"}}, Player: "testUser"})

	g := NewWithT(t)
	g.Expect(*eventSender).To(Equal(mockEventSender{receivedCode: "testCode", receivedMessage: toJsonBytes(clientEvent{PlayerJoined: &playerJoined{Name: "testUser"}})}))
}
