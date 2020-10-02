package messagebus

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testConsumer struct {
	receivedMessage string
}

func (t *testConsumer) Consume(m Message) {
	t.receivedMessage += m.(testMessage).message
}

type testMessage struct {
	message string
}

func (t testMessage) GetPartyCode() string {
	return ""
}

func (t testMessage) Type() Type {
	return Type("")
}

func Test_DispatchMessage(t *testing.T) {
	mb := NewMessageBus()

	testConsumer1 := &testConsumer{}
	testConsumer2 := &testConsumer{}
	testConsumer3 := &testConsumer{}
	mb.SubscribeConsumer(testConsumer1)
	mb.SubscribeConsumer(testConsumer2)
	mb.SubscribeConsumer(testConsumer3)

	mb.Dispatch(testMessage{"m1"})
	mb.Dispatch(testMessage{"m2"})
	mb.Dispatch(testMessage{"m3"})
	mb.Close()

	g := NewWithT(t)
	g.Expect(testConsumer1.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer2.receivedMessage).To(Equal("m1m2m3"))
	g.Expect(testConsumer3.receivedMessage).To(Equal("m1m2m3"))
}

func Test_GetPartyCode(t *testing.T) {
	m := Party{Code: "testCode"}

	g := NewWithT(t)
	g.Expect(m.GetPartyCode()).To(Equal("testCode"))
}

func Test_TypeCommand(t *testing.T) {
	m := Command{Party{Code: "testCode"}}

	g := NewWithT(t)
	g.Expect(m.Type()).To(Equal(CommandMessage))
}

func Test_TypeEvent(t *testing.T) {
	m := Event{Party{Code: "testCode"}}

	g := NewWithT(t)
	g.Expect(m.Type()).To(Equal(EventMessage))
}
