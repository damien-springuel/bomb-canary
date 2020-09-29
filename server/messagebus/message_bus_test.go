package messagebus

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testConsumer struct {
	receivedMessage string
}

func (t *testConsumer) consume(m Message) {
	t.receivedMessage += m.(testMessage).message
}

type testMessage struct {
	message string
}

func (t testMessage) GetPartyCode() string {
	return ""
}

func Test_DispatchMessage(t *testing.T) {
	mb := NewMessageBus()

	testConsumer1 := &testConsumer{}
	testConsumer2 := &testConsumer{}
	testConsumer3 := &testConsumer{}
	mb.subscribeConsumer(testConsumer1)
	mb.subscribeConsumer(testConsumer2)
	mb.subscribeConsumer(testConsumer3)

	mb.dispatchMessage(testMessage{"m1"})
	mb.dispatchMessage(testMessage{"m2"})
	mb.dispatchMessage(testMessage{"m3"})
	mb.close()

	g := NewWithT(t)
	g.Expect(testConsumer1.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer2.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer3.receivedMessage).To(Equal("m1m2m3"))
}
