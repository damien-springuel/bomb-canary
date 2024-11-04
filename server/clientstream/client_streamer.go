package clientstream

import (
	"sync"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type name string

type messageDispatcher interface {
	Dispatch(m messagebus.Message)
}

type clientStreamer struct {
	mut               *sync.RWMutex
	clientOutByName   map[name]chan []byte
	messageDispatcher messageDispatcher
}

func NewClientsStreamer(messageDispatcher messageDispatcher) clientStreamer {
	return clientStreamer{
		mut:               &sync.RWMutex{},
		clientOutByName:   make(map[name]chan []byte),
		messageDispatcher: messageDispatcher,
	}
}

func (c clientStreamer) Add(playerName string) (chan []byte, func()) {
	c.mut.Lock()
	defer c.mut.Unlock()

	clientOut := make(chan []byte)
	_, exists := c.clientOutByName[name(playerName)]
	if !exists {
		c.clientOutByName[name(playerName)] = clientOut
	}

	c.dispatchConnectedMessage(playerName)

	return clientOut, func() {
		c.remove(playerName)
	}
}

func (c clientStreamer) remove(playerName string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	client, exists := c.clientOutByName[name(playerName)]
	if exists {
		close(client)
		c.dispatchDisconnectedMessage(playerName)
	}
	if exists {
	}
}

func (c clientStreamer) Send(message []byte) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	for _, out := range c.clientOutByName {
		out <- []byte(message)
	}
}

func (c clientStreamer) SendToPlayer(playerName string, message []byte) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	client, exists := c.clientOutByName[name(playerName)]
	if exists {
		client <- message
	}
}

func (c clientStreamer) SendToAllButPlayer(playerName string, message []byte) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	for n, out := range c.clientOutByName {
		if n != name(playerName) {
			out <- []byte(message)
		}
	}
}

func (c clientStreamer) dispatchConnectedMessage(name string) {
	c.messageDispatcher.Dispatch(messagebus.PlayerConnected{Player: name})
}

func (c clientStreamer) dispatchDisconnectedMessage(name string) {
	c.messageDispatcher.Dispatch(messagebus.PlayerDisconnected{Player: name})
}
