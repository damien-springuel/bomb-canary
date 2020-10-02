package clientstream

import (
	"sync"
)

type code string
type name string

type clientStreamer struct {
	mut                   *sync.RWMutex
	clientOutByNameByCode map[code]map[name]chan []byte
}

func NewClientsStreamer() clientStreamer {
	return clientStreamer{
		mut:                   &sync.RWMutex{},
		clientOutByNameByCode: make(map[code]map[name]chan []byte),
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
	return clientOut, func() {
		c.remove(partyCode, playerName)
	}
}

func (c clientStreamer) remove(partyCode, playerName string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	clients := c.clientOutByNameByCode[code(partyCode)]
	out := clients[name(playerName)]
	close(out)
	delete(clients, name(playerName))
	if len(clients) == 0 {
		delete(c.clientOutByNameByCode, code(partyCode))
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
