package clientstream

import (
	"fmt"
	"sync"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type code string
type name string

type client struct {
	out chan string
}

type clientStreamer struct {
	mut                 *sync.RWMutex
	clientsByNameByCode map[code]map[name]client
}

func NewClientsStreamer() clientStreamer {
	return clientStreamer{
		mut:                 &sync.RWMutex{},
		clientsByNameByCode: make(map[code]map[name]client),
	}
}

func (c clientStreamer) Add(partyCode, playerName string) chan string {
	clientOut := make(chan string)
	c.mut.Lock()
	defer c.mut.Unlock()
	clients, exists := c.clientsByNameByCode[code(partyCode)]
	if exists {
		clients[name(playerName)] = client{out: clientOut}
	} else {
		clients = map[name]client{
			name(playerName): {out: clientOut},
		}
	}
	c.clientsByNameByCode[code(partyCode)] = clients
	return clientOut
}

func (c clientStreamer) sendMessageToParty(partyCode, message string) {
	c.mut.RLock()
	defer c.mut.RUnlock()
	clients, exists := c.clientsByNameByCode[code(partyCode)]
	if exists {
		for _, client := range clients {
			client.out <- message
		}
	}
}

func (c clientStreamer) Consume(m messagebus.Message) {
	code := m.GetPartyCode()
	c.sendMessageToParty(code, fmt.Sprintf("%#v", m))
}
