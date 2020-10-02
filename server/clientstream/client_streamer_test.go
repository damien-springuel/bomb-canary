package clientstream

import (
	"testing"

	. "github.com/onsi/gomega"
)

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

	streamer := NewClientsStreamer()
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

	streamer := NewClientsStreamer()
	testOut := make(chan [][]byte)
	done := createAndPumpOut(streamer, "testCode", "p1", testOut)

	streamer.Send("other", []byte("message1"))

	done()
	actualMessages := <-testOut

	g := NewWithT(t)
	g.Expect(actualMessages).To(BeEmpty())
}

func Test_SendMultiplePlayersInParty(t *testing.T) {

	streamer := NewClientsStreamer()

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

	streamer := NewClientsStreamer()

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
