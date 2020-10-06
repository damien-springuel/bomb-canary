package clientstream

import (
	"sync"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type code string
type name string

type messageDispatcher interface {
	Dispatch(m messagebus.Message)
}

type clientStreamer struct {
	mut                   *sync.RWMutex
	clientOutByNameByCode map[code]map[name]chan []byte
	messageDispatcher     messageDispatcher
}

func NewClientsStreamer(messageDispatcher messageDispatcher) clientStreamer {
	return clientStreamer{
		mut:                   &sync.RWMutex{},
		clientOutByNameByCode: make(map[code]map[name]chan []byte),
		messageDispatcher:     messageDispatcher,
	}
}

func (c clientStreamer) Add(partyCode, playerName string) (chan []byte, func()) {
	c.mut.Lock()
	defer c.mut.Unlock()

	clientOut := make(chan []byte)
	clients, exists := c.clientOutByNameByCode[code(partyCode)]
	if exists {
		clients[name(playerName)] = clientOut
	} else {
		clients = map[name]chan []byte{
			name(playerName): clientOut,
		}
	}
	c.clientOutByNameByCode[code(partyCode)] = clients

	c.dispatchConnectedMessage(partyCode, playerName)

	return clientOut, func() {
		c.remove(partyCode, playerName)
	}
}

func (c clientStreamer) remove(partyCode, playerName string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	clients, exists := c.clientOutByNameByCode[code(partyCode)]
	if exists {
		out, exists := clients[name(playerName)]
		if exists {
			close(out)
			delete(clients, name(playerName))
			if len(clients) == 0 {
				delete(c.clientOutByNameByCode, code(partyCode))
			}
			c.dispatchDisconnectedMessage(partyCode, playerName)
		}
	}
}

func (c clientStreamer) Send(partyCode string, message []byte) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	clients, exists := c.clientOutByNameByCode[code(partyCode)]
	if exists {
		for _, out := range clients {
			out <- []byte(message)
		}
	}
}

func (c clientStreamer) SendToPlayer(partyCode, playerName string, message []byte) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	clients, exists := c.clientOutByNameByCode[code(partyCode)]
	if exists {
		out, exists := clients[name(playerName)]
		if exists {
			out <- message
		}
	}
}

func (c clientStreamer) SendToAllButPlayer(partyCode, playerName string, message []byte) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	clients, exists := c.clientOutByNameByCode[code(partyCode)]
	if exists {
		for n, out := range clients {
			if n != name(playerName) {
				out <- []byte(message)
			}
		}
	}
}

func (c clientStreamer) dispatchConnectedMessage(code, name string) {
	c.messageDispatcher.Dispatch(messagebus.PlayerConnected{Event: messagebus.Event{Party: messagebus.Party{Code: code}}, Player: name})
}

func (c clientStreamer) dispatchDisconnectedMessage(code, name string) {
	c.messageDispatcher.Dispatch(messagebus.PlayerDisconnected{Event: messagebus.Event{Party: messagebus.Party{Code: code}}, Player: name})
}
