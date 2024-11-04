package clientstream

import (
	"encoding/json"
	"testing"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
	. "github.com/onsi/gomega"
)

var expectedReplayEnded, _ = json.Marshal(clientEvent{EventsReplayEnded: &eventsReplayEnded{}})

func Test_Replayer_NothingToReplay(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.Consume(messagebus.PlayerConnected{Player: "testName"})

	expectedReplayStarted, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "testName"}})

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStarted,
		expectedReplayEnded,
	}))
}

func Test_Replayer_Send(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.Send([]byte("m1"))
	replayer.Send([]byte("m2"))
	replayer.Send([]byte("m3"))
	replayer.Send([]byte("m4"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Player: "testName"})

	expectedReplayStarted, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "testName"}})
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStarted,
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
		expectedReplayEnded,
	}))
}

func Test_Replayer_SendToPlayer(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.SendToPlayer("p1", []byte("m1"))
	replayer.SendToPlayer("p2", []byte("m2"))
	replayer.SendToPlayer("p1", []byte("m3"))
	replayer.SendToPlayer("p1", []byte("m4"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Player: "p1"})

	expectedReplayStartedP1, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "p1"}})
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStartedP1,
		[]byte("m1"),
		[]byte("m3"),
		[]byte("m4"),
		expectedReplayEnded,
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Player: "p2"})

	expectedReplayStartedP2, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "p2"}})
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStartedP2,
		[]byte("m2"),
		expectedReplayEnded,
	}))
}

func Test_Replayer_SendToAllButPlayer(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.SendToAllButPlayer("p1", []byte("m1"))
	replayer.SendToAllButPlayer("p2", []byte("m2"))
	replayer.SendToAllButPlayer("p1", []byte("m3"))
	replayer.SendToAllButPlayer("p1", []byte("m4"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Player: "p1"})

	expectedReplayStartedP1, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "p1"}})
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStartedP1,
		[]byte("m2"),
		expectedReplayEnded,
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Player: "p2"})

	expectedReplayStartedP2, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "p2"}})
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStartedP2,
		[]byte("m1"),
		[]byte("m3"),
		[]byte("m4"),
		expectedReplayEnded,
	}))
}

func Test_Replayer_MixedCase(t *testing.T) {
	mockEventSender := &mockEventSender{shouldTrackAll: true}
	replayer := NewEventReplayer(mockEventSender)
	replayer.Send([]byte("m1"))
	replayer.Send([]byte("m2"))
	replayer.SendToPlayer("p1", []byte("m3"))
	replayer.SendToAllButPlayer("p1", []byte("m4"))
	replayer.SendToAllButPlayer("p2", []byte("m5"))
	replayer.SendToPlayer("p2", []byte("m6"))
	replayer.Send([]byte("m7"))

	g := NewWithT(t)
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m4"),
		[]byte("m5"),
		[]byte("m6"),
		[]byte("m7"),
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Player: "p1"})

	expectedReplayStartedP1, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "p1"}})
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStartedP1,
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m3"),
		[]byte("m5"),
		[]byte("m7"),
		expectedReplayEnded,
	}))
	mockEventSender.clearAllReceivedMessages()

	replayer.Consume(messagebus.PlayerConnected{Player: "p2"})

	expectedReplayStartedP2, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: "p2"}})
	g.Expect(mockEventSender.allReceivedMessages).To(Equal([][]byte{
		expectedReplayStartedP2,
		[]byte("m1"),
		[]byte("m2"),
		[]byte("m4"),
		[]byte("m6"),
		[]byte("m7"),
		expectedReplayEnded,
	}))
}
