package sessions

import (
	"errors"
	"sync"
)

type uuidCreator interface {
	Create() string
}

type sessions struct {
	uuidCreator     uuidCreator
	mut             *sync.RWMutex
	nameBySessionId map[string]string
}

func New(uuidCreator uuidCreator) sessions {
	return sessions{
		uuidCreator:     uuidCreator,
		nameBySessionId: make(map[string]string),
		mut:             &sync.RWMutex{},
	}
}

func (s sessions) Create(name string) string {
	s.mut.Lock()
	defer s.mut.Unlock()

	uuid := s.uuidCreator.Create()

	s.nameBySessionId[uuid] = name

	return uuid
}

func (s sessions) Get(session string) (name string, err error) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	playerName, exists := s.nameBySessionId[session]
	if !exists {
		return "", errors.New("session doesn't exist")
	}

	return playerName, nil
}
