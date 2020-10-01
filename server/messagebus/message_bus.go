package messagebus

type Message interface {
	GetPartyCode() string
}

type consumer interface {
	Consume(m Message)
}

type bufferedConsumer struct {
	consumer consumer
	in       chan Message
	done     chan bool
}

func newBufferedConsumer(consumer consumer) bufferedConsumer {
	bc := bufferedConsumer{
		consumer: consumer,
		in:       make(chan Message, 10),
		done:     make(chan bool),
	}
	go func() {
		for incomingMessage := range bc.in {
			bc.consumer.Consume(incomingMessage)
		}
		bc.done <- true
	}()

	return bc
}

func (b bufferedConsumer) consume(m Message) {
	b.in <- m
}

func (b bufferedConsumer) stop() {
	close(b.in)
	<-b.done
}

type messageBus struct {
	consumers []bufferedConsumer
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

func (mb *messageBus) Dispatch(m Message) {
	mb.in <- m
}

func (m *messageBus) SubscribeConsumer(consumer consumer) {
	m.consumers = append(m.consumers, newBufferedConsumer(consumer))
}

func (m *messageBus) close() {
	close(m.in)
	<-m.done
	for _, bc := range m.consumers {
		bc.stop()
	}
}
