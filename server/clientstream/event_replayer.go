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

	mut            *sync.RWMutex
	messagesByCode map[code][]replayMessage
}

func NewEventReplayer(eventSender eventSender) eventReplayer {
	return eventReplayer{
		eventSender:    eventSender,
		mut:            &sync.RWMutex{},
		messagesByCode: make(map[code][]replayMessage),
	}
}

func (e eventReplayer) Consume(m messagebus.Message) {
	connectEvent, ok := m.(messagebus.PlayerConnected)
	if ok {
		e.mut.RLock()
		defer e.mut.RUnlock()

		replayStartedMessage, _ := json.Marshal(clientEvent{EventsReplayStarted: &eventsReplayStarted{Player: connectEvent.Player}})
		e.eventSender.SendToPlayer(connectEvent.Code, connectEvent.Player, replayStartedMessage)

		e.sendReplayableMessages(connectEvent.Code, connectEvent.Player)

		replayEndedMessage, _ := json.Marshal(clientEvent{EventsReplayEnded: &eventsReplayEnded{}})
		e.eventSender.SendToPlayer(connectEvent.Code, connectEvent.Player, replayEndedMessage)
	}
}

func (e eventReplayer) sendReplayableMessages(partyCode, playerName string) {
	messages, exists := e.messagesByCode[code(partyCode)]
	if exists {
		for _, replayMessage := range messages {
			if replayMessage.replayType == All ||
				(replayMessage.replayType == Player && replayMessage.name == playerName) ||
				(replayMessage.replayType == AllButPlayer && replayMessage.name != playerName) {
				e.eventSender.SendToPlayer(partyCode, playerName, replayMessage.message)
			}
		}
	}
}

func (e eventReplayer) recordMessage(partyCode string, replayMessage replayMessage) {
	e.mut.Lock()
	defer e.mut.Unlock()
	messages := e.messagesByCode[code(partyCode)]
	messages = append(messages, replayMessage)
	e.messagesByCode[code(partyCode)] = messages
}

func (e eventReplayer) Send(partyCode string, message []byte) {
	e.recordMessage(partyCode, replayMessage{replayType: All, message: message})
	e.eventSender.Send(partyCode, message)
}

func (e eventReplayer) SendToPlayer(partyCode, playerName string, message []byte) {
	e.recordMessage(partyCode, replayMessage{replayType: Player, name: playerName, message: message})
	e.eventSender.SendToPlayer(partyCode, playerName, message)
}

func (e eventReplayer) SendToAllButPlayer(partyCode, playerName string, message []byte) {
	e.recordMessage(partyCode, replayMessage{replayType: AllButPlayer, name: playerName, message: message})
	e.eventSender.SendToAllButPlayer(partyCode, playerName, message)
}
