package clientstream

import (
	"testing"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

type mockMessageDispatcher struct {
	receivedMessages []messagebus.Message
}

func (d *mockMessageDispatcher) Dispatch(m messagebus.Message) {
	d.receivedMessages = append(d.receivedMessages, m)
}

func createAndPumpOut(streamer clientStreamer, code, name string, done chan [][]byte) func() {
	out, closer := streamer.Add(code, name)

	actualMessages := [][]byte{}
	go func() {
		for nextMessage := range out {
			actualMessages = append(actualMessages, nextMessage)
		}
		done <- actualMessages
		close(done)
	}()

	return closer
}

func Test_Send(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})
	testOut := make(chan [][]byte)
	done := createAndPumpOut(streamer, "testCode", "p1", testOut)

	streamer.Send("testCode", []byte("message1"))
	streamer.Send("testCode", []byte("message2"))

	done()
	actualMessages := <-testOut

	g := NewWithT(t)
	g.Expect(actualMessages).To(Equal([][]byte{[]byte("message1"), []byte("message2")}))
}

func Test_SendOnAnotherCode(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})
	testOut := make(chan [][]byte)
	done := createAndPumpOut(streamer, "testCode", "p1", testOut)

	streamer.Send("other", []byte("message1"))

	done()
	actualMessages := <-testOut

	g := NewWithT(t)
	g.Expect(actualMessages).To(BeEmpty())
}

func Test_SendMultiplePlayersInParty(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})

	testOut1 := make(chan [][]byte)
	done1 := createAndPumpOut(streamer, "testCode", "p1", testOut1)
	testOut2 := make(chan [][]byte)
	done2 := createAndPumpOut(streamer, "testCode", "p2", testOut2)

	streamer.Send("testCode", []byte("m1"))
	streamer.Send("testCode", []byte("m2"))

	done1()
	done2()
	actualMessages1 := <-testOut1
	actualMessages2 := <-testOut2

	g := NewWithT(t)
	g.Expect(actualMessages1).To(Equal([][]byte{[]byte("m1"), []byte("m2")}))
	g.Expect(actualMessages2).To(Equal([][]byte{[]byte("m1"), []byte("m2")}))
}

func Test_SendMultiplePlayersInMultipleParty(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})

	testOut1 := make(chan [][]byte)
	done1 := createAndPumpOut(streamer, "c1", "1p1", testOut1)
	testOut2 := make(chan [][]byte)
	done2 := createAndPumpOut(streamer, "c1", "1p2", testOut2)

	testOut3 := make(chan [][]byte)
	done3 := createAndPumpOut(streamer, "c2", "2p1", testOut3)
	testOut4 := make(chan [][]byte)
	done4 := createAndPumpOut(streamer, "c2", "2p2", testOut4)

	streamer.Send("c1", []byte("m1"))
	streamer.Send("c2", []byte("m2"))
	streamer.Send("c1", []byte("m3"))
	streamer.Send("c1", []byte("m4"))
	streamer.Send("c2", []byte("m5"))
	streamer.Send("c2", []byte("m6"))
	streamer.Send("c1", []byte("m7"))
	streamer.Send("c2", []byte("m8"))

	done1()
	done2()
	done3()
	done4()
	actualMessages1 := <-testOut1
	actualMessages2 := <-testOut2
	actualMessages3 := <-testOut3
	actualMessages4 := <-testOut4

	g := NewWithT(t)
	g.Expect(actualMessages1).To(Equal([][]byte{[]byte("m1"), []byte("m3"), []byte("m4"), []byte("m7")}))
	g.Expect(actualMessages2).To(Equal([][]byte{[]byte("m1"), []byte("m3"), []byte("m4"), []byte("m7")}))
	g.Expect(actualMessages3).To(Equal([][]byte{[]byte("m2"), []byte("m5"), []byte("m6"), []byte("m8")}))
	g.Expect(actualMessages4).To(Equal([][]byte{[]byte("m2"), []byte("m5"), []byte("m6"), []byte("m8")}))
}

func Test_SendToParticularPlayer(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})

	testOut1 := make(chan [][]byte)
	done1 := createAndPumpOut(streamer, "c1", "1p1", testOut1)
	testOut2 := make(chan [][]byte)
	done2 := createAndPumpOut(streamer, "c1", "1p2", testOut2)
	testOut3 := make(chan [][]byte)
	done3 := createAndPumpOut(streamer, "c2", "2p1", testOut3)
	testOut4 := make(chan [][]byte)
	done4 := createAndPumpOut(streamer, "c2", "2p2", testOut4)

	streamer.Send("c1", []byte("m1"))
	streamer.SendToPlayer("c1", "1p1", []byte("m2"))
	streamer.SendToPlayer("c1", "2p1", []byte("m3"))

	done1()
	done2()
	done3()
	done4()
	actualMessages1 := <-testOut1
	actualMessages2 := <-testOut2
	actualMessages3 := <-testOut3
	actualMessages4 := <-testOut4

	g := NewWithT(t)
	g.Expect(actualMessages1).To(Equal([][]byte{[]byte("m1"), []byte("m2")}))
	g.Expect(actualMessages2).To(Equal([][]byte{[]byte("m1")}))
	g.Expect(actualMessages3).To(BeEmpty())
	g.Expect(actualMessages4).To(BeEmpty())
}

func Test_SendToAllButPlayer(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})

	testOut1 := make(chan [][]byte)
	done1 := createAndPumpOut(streamer, "c1", "p1", testOut1)
	testOut2 := make(chan [][]byte)
	done2 := createAndPumpOut(streamer, "c1", "p2", testOut2)
	testOut3 := make(chan [][]byte)
	done3 := createAndPumpOut(streamer, "c1", "p3", testOut3)
	testOut4 := make(chan [][]byte)
	done4 := createAndPumpOut(streamer, "c1", "p4", testOut4)

	streamer.Send("c1", []byte("m1"))
	streamer.SendToPlayer("c1", "p1", []byte("m2"))
	streamer.SendToAllButPlayer("c1", "p1", []byte("m3"))

	done1()
	done2()
	done3()
	done4()
	actualMessages1 := <-testOut1
	actualMessages2 := <-testOut2
	actualMessages3 := <-testOut3
	actualMessages4 := <-testOut4

	g := NewWithT(t)
	g.Expect(actualMessages1).To(Equal([][]byte{[]byte("m1"), []byte("m2")}))
	g.Expect(actualMessages2).To(Equal([][]byte{[]byte("m1"), []byte("m3")}))
	g.Expect(actualMessages3).To(Equal([][]byte{[]byte("m1"), []byte("m3")}))
	g.Expect(actualMessages4).To(Equal([][]byte{[]byte("m1"), []byte("m3")}))
}

func Test_AddAndRemoveDispatchPlayerConnectedDisconnectedMessage(t *testing.T) {
	dispatcher := &mockMessageDispatcher{}
	streamer := NewClientsStreamer(dispatcher)
	testOut := make(chan [][]byte)
	done := createAndPumpOut(streamer, "testCode", "p1", testOut)

	streamer.Send("testCode", []byte("message1"))
	streamer.Send("testCode", []byte("message2"))

	done()
	<-testOut

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessages).To(Equal([]messagebus.Message{
		messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p1"},
		messagebus.PlayerDisconnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p1"},
	}))
}
