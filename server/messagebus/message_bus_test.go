package messagebus

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testConsumer struct {
	receivedMessage string
	name            string
}

func (t *testConsumer) consume(m Message) {
	t.receivedMessage += m.(string)
}

func Test_DispatchMessage(t *testing.T) {
	mb := NewMessageBus()

	testConsumer1 := &testConsumer{name: "1"}
	testConsumer2 := &testConsumer{name: "2"}
	testConsumer3 := &testConsumer{name: "3"}
	mb.SubscribeConsumer(testConsumer1)
	mb.SubscribeConsumer(testConsumer2)
	mb.SubscribeConsumer(testConsumer3)

	mb.dispatchMessage("m1")
	mb.dispatchMessage("m2")
	mb.dispatchMessage("m3")
	mb.stop()

	g := NewWithT(t)
	g.Expect(testConsumer1.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer2.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer3.receivedMessage).To(Equal("m1m2m3"))
}
