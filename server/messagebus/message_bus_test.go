package messagebus

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testConsumer struct {
	receivedMessage Message
}

func (t *testConsumer) consume(m Message) {
	t.receivedMessage = m
}

func Test_DispatchMessage(t *testing.T) {
	mb := NewMessageBus()
	testConsumer1 := &testConsumer{receivedMessage: "old 1"}
	testConsumer2 := &testConsumer{receivedMessage: "old 2"}
	testConsumer3 := &testConsumer{receivedMessage: "old 3"}
	mb.SubscribeConsumer(testConsumer1)
	mb.SubscribeConsumer(testConsumer2)
	mb.SubscribeConsumer(testConsumer3)

	mb.dispatchMessage("new message")

	g := NewWithT(t)
	g.Expect(testConsumer1.receivedMessage).To(Equal("new message"))
	g.Expect(testConsumer2.receivedMessage).To(Equal("new message"))
	g.Expect(testConsumer3.receivedMessage).To(Equal("new message"))
}
