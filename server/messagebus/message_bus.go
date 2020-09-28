package messagebus

type Message interface{}

type consumer interface {
	consume(m Message)
}

type messageBus struct {
	consumers []consumer
}

func NewMessageBus() messageBus {
	return messageBus{}
}

func (mb *messageBus) dispatchMessage(m Message) {
	for _, c := range mb.consumers {
		c.consume(m)
	}
}

func (m *messageBus) SubscribeConsumer(consumer consumer) {
	m.consumers = append(m.consumers, consumer)
}
