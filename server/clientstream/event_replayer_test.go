package clientstream

import (
	"encoding/json"
	"testing"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

var expectedReplayEnded, _ = json.Marshal(struct{ ReplayEnded struct{} }{ReplayEnded: struct{}{}})

func Test_Replayer_NothingToReplay(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "testName"})

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayEnded,
	}))
}

func Test_Replayer_Send(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.Send("testCode", []byte("m1"))
	replayer.Send("testCode", []byte("m2"))
	replayer.Send("other", []byte("m3"))
	replayer.Send("testCode", []byte("m4"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "testName"})

	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m4"),
		expectedReplayEnded,
	}))
}

func Test_Replayer_SendToPlayer(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.SendToPlayer("testCode", "p1", []byte("m1"))
	replayer.SendToPlayer("testCode", "p2", []byte("m2"))
	replayer.SendToPlayer("otherCode", "p1", []byte("m3"))
	replayer.SendToPlayer("testCode", "p1", []byte("m4"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p1"})

	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m4"),
		expectedReplayEnded,
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p2"})

	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m2"),
		expectedReplayEnded,
	}))
}

func Test_Replayer_SendToAllButPlayer(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.SendToAllButPlayer("testCode", "p1", []byte("m1"))
	replayer.SendToAllButPlayer("testCode", "p2", []byte("m2"))
	replayer.SendToAllButPlayer("otherCode", "p1", []byte("m3"))
	replayer.SendToAllButPlayer("testCode", "p1", []byte("m4"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p1"})

	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m2"),
		expectedReplayEnded,
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p2"})

	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m4"),
		expectedReplayEnded,
	}))
}

func Test_Replayer_MixedCase(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.Send("testCode", []byte("m1"))
	replayer.Send("otherCode", []byte("m2"))
	replayer.SendToPlayer("testCode", "p1", []byte("m3"))
	replayer.SendToAllButPlayer("testCode", "p2", []byte("m4"))
	replayer.SendToPlayer("otherCode", "p1", []byte("m5"))
	replayer.Send("testCode", []byte("m6"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
		[]byte("m5"),
		[]byte("m6"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p1"})

	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m3"),
		[]byte("m4"),
		[]byte("m6"),
		expectedReplayEnded,
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: "testCode"}}, Player: "p2"})

	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m6"),
		expectedReplayEnded,
	}))
}
