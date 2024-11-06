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

func createAndPumpOut(streamer clientStreamer, name string, done chan [][]byte) func() {
	out, closer := streamer.Add(name)

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
	done := createAndPumpOut(streamer, "p1", testOut)

	streamer.Send([]byte("message1"))
	streamer.Send([]byte("message2"))

	done()
	actualMessages := <-testOut

	g := NewWithT(t)
	g.Expect(actualMessages).To(Equal([][]byte{[]byte("message1"), []byte("message2")}))
}

func Test_SendMultiplePlayersInParty(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})

	testOut1 := make(chan [][]byte)
	done1 := createAndPumpOut(streamer, "p1", testOut1)
	testOut2 := make(chan [][]byte)
	done2 := createAndPumpOut(streamer, "p2", testOut2)

	streamer.Send([]byte("m1"))
	streamer.Send([]byte("m2"))

	done1()
	done2()
	actualMessages1 := <-testOut1
	actualMessages2 := <-testOut2

	g := NewWithT(t)
	g.Expect(actualMessages1).To(Equal([][]byte{[]byte("m1"), []byte("m2")}))
	g.Expect(actualMessages2).To(Equal([][]byte{[]byte("m1"), []byte("m2")}))
}

func Test_SendToAllButPlayer(t *testing.T) {
	streamer := NewClientsStreamer(&mockMessageDispatcher{})

	testOut1 := make(chan [][]byte)
	done1 := createAndPumpOut(streamer, "p1", testOut1)
	testOut2 := make(chan [][]byte)
	done2 := createAndPumpOut(streamer, "p2", testOut2)
	testOut3 := make(chan [][]byte)
	done3 := createAndPumpOut(streamer, "p3", testOut3)
	testOut4 := make(chan [][]byte)
	done4 := createAndPumpOut(streamer, "p4", testOut4)

	streamer.Send([]byte("m1"))
	streamer.SendToPlayer("p1", []byte("m2"))
	streamer.SendToAllButPlayer("p1", []byte("m3"))

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
	done := createAndPumpOut(streamer, "p1", testOut)

	streamer.Send([]byte("message1"))
	streamer.Send([]byte("message2"))

	done()
	<-testOut

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessages).To(Equal([]messagebus.Message{
		messagebus.PlayerConnected{Player: "p1"},
		messagebus.PlayerDisconnected{Player: "p1"},
	}))
}
func Test_AddAndRemove_AndReconnectASecondTime(t *testing.T) {
	dispatcher := &mockMessageDispatcher{}
	streamer := NewClientsStreamer(dispatcher)
	testOut := make(chan [][]byte)
	done := createAndPumpOut(streamer, "p1", testOut)

	streamer.Send([]byte("message1"))
	streamer.Send([]byte("message2"))

	done()
	<-testOut

	testOut = make(chan [][]byte)
	done = createAndPumpOut(streamer, "p1", testOut)
	streamer.Send([]byte("message3"))
	streamer.Send([]byte("message4"))

	done()
	<-testOut

	g := NewWithT(t)
	g.Expect(dispatcher.receivedMessages).To(Equal([]messagebus.Message{
		messagebus.PlayerConnected{Player: "p1"},
		messagebus.PlayerDisconnected{Player: "p1"},
		messagebus.PlayerConnected{Player: "p1"},
		messagebus.PlayerDisconnected{Player: "p1"},
	}))
}
