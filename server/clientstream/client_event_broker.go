package clientstream

import (
	"fmt"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type eventSender interface {
	Send(code string, message []byte)
}

type clientEventBroker struct {
	eventSender eventSender
}

func NewClientEventBroker(eventSender eventSender) clientEventBroker {
	return clientEventBroker{
		eventSender: eventSender,
	}
}

func (c clientEventBroker) Consume(m messagebus.Message) {
	code := m.GetPartyCode()
	c.eventSender.Send(code, []byte(fmt.Sprintf("%#v", m)))
}
