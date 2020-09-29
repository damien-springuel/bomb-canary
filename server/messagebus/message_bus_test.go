package messagebus

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testConsumer struct {
	receivedMessage string
}

func (t *testConsumer) consume(m Message) {
	t.receivedMessage += m.(string)
}

func Test_DispatchMessage(t *testing.T) {
	mb := NewMessageBus()

	testConsumer1 := &testConsumer{}
	testConsumer2 := &testConsumer{}
	testConsumer3 := &testConsumer{}
	mb.SubscribeConsumer(testConsumer1)
	mb.SubscribeConsumer(testConsumer2)
	mb.SubscribeConsumer(testConsumer3)

	mb.dispatchMessage("m1")
	mb.dispatchMessage("m2")
	mb.dispatchMessage("m3")
	mb.close()

	g := NewWithT(t)
	g.Expect(testConsumer1.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer2.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer3.receivedMessage).To(Equal("m1m2m3"))
}
