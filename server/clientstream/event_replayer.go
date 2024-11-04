package clientstream

import (
	"encoding/json"
	"sync"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type replayType string

const (
	All          replayType = "all"
	Player       replayType = "player"
	AllButPlayer replayType = "allButPlayer"
)

type replayMessage struct {
	replayType replayType
	name       string
	message    []byte
}

type eventReplayer struct {
	eventSender eventSender
	mut         *sync.RWMutex
	messages    []replayMessage
}

func NewEventReplayer(eventSender eventSender) *eventReplayer {
	return &eventReplayer{
		eventSender: eventSender,
		mut:         &sync.RWMutex{},
		messages:    make([]replayMessage, 0),
	}
}

func (e *eventReplayer) Consume(m messagebus.Message) {
	connectEvent, ok := m.(messagebus.PlayerConnected)
	if ok {
		e.mut.RLock()
		defer e.mut.RUnlock()

		replayStartedMessage, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: connectEvent.Player}})
		e.eventSender.SendToPlayer(connectEvent.Player, replayStartedMessage)

		e.sendReplayableMessages(connectEvent.Player)

		replayEndedMessage, _ := json.Marshal(clientEvent{EventsReplayEnded: &eventsReplayEnded{}})
		e.eventSender.SendToPlayer(connectEvent.Player, replayEndedMessage)
	}
}

func (e *eventReplayer) sendReplayableMessages(playerName string) {
	for _, replayMessage := range e.messages {
		if replayMessage.replayType == All ||
			(replayMessage.replayType == Player && replayMessage.name == playerName) ||
			(replayMessage.replayType == AllButPlayer && replayMessage.name != playerName) {
			e.eventSender.SendToPlayer(playerName, replayMessage.message)
		}
	}
}

func (e *eventReplayer) recordMessage(replayMessage replayMessage) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.messages = append(e.messages, replayMessage)
}

func (e *eventReplayer) Send(message []byte) {
	e.recordMessage(replayMessage{replayType: All, message: message})
	e.eventSender.Send(message)
}

func (e *eventReplayer) SendToPlayer(playerName string, message []byte) {
	e.recordMessage(replayMessage{replayType: Player, name: playerName, message: message})
	e.eventSender.SendToPlayer(playerName, message)
}

func (e *eventReplayer) SendToAllButPlayer(playerName string, message []byte) {
	e.recordMessage(replayMessage{replayType: AllButPlayer, name: playerName, message: message})
	e.eventSender.SendToAllButPlayer(playerName, message)
}
