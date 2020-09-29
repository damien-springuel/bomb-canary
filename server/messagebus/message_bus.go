package messagebus

type Message interface{}

type consumer interface {
	consume(m Message)
}

type messageBus struct {
	consumers []consumer
	in        chan Message
	done      chan bool
}

func NewMessageBus() *messageBus {
	mb := messageBus{in: make(chan Message, 10), done: make(chan bool)}

	go func() {
		for incomingMessage := range mb.in {
			for _, c := range mb.consumers {
				c.consume(incomingMessage)
			}
		}
		mb.done <- true
	}()

	return &mb
}

func (mb *messageBus) dispatchMessage(m Message) {
	mb.in <- m
}

func (m *messageBus) SubscribeConsumer(consumer consumer) {
	m.consumers = append(m.consumers, consumer)
}

func (m *messageBus) stop() {
	close(m.in)
	<-m.done
}
