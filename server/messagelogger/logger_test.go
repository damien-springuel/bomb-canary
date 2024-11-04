package messagelogger

import (
	"testing"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type mockPrinter struct {
	receivedCommand messagebus.Message
	receivedEvent   messagebus.Message
}

func (m *mockPrinter) PrintCommand(message messagebus.Message) {
	m.receivedCommand = message
}

func (m *mockPrinter) PrintEvent(message messagebus.Message) {
	m.receivedEvent = message
}

func Test_Consume_Command(t *testing.T) {
	printer := &mockPrinter{}
	l := New(printer)
	command := messagebus.JoinParty{Player: "name"}
	l.Consume(command)

	g := NewWithT(t)
	g.Expect(printer.receivedCommand).To(Equal(messagebus.JoinParty{Player: "name"}))
	g.Expect(printer.receivedEvent).To(BeNil())
}

func Test_Consume_Event(t *testing.T) {
	printer := &mockPrinter{}
	l := New(printer)
	command := messagebus.PlayerJoined{Player: "name"}
	l.Consume(command)

	g := NewWithT(t)
	g.Expect(printer.receivedEvent).To(Equal(messagebus.PlayerJoined{Player: "name"}))
	g.Expect(printer.receivedCommand).To(BeNil())
}
