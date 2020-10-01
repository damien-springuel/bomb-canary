package sessions

import "sync"

type playerSession struct {
	partyCode string
	name      string
}

type uuidCreator interface {
	Create() string
}

type sessions struct {
	uuidCreator uuidCreator
	mut         *sync.RWMutex
	sessionById map[string]playerSession
}

func New(uuidCreator uuidCreator) sessions {
	return sessions{
		uuidCreator: uuidCreator,
		sessionById: make(map[string]playerSession),
		mut:         &sync.RWMutex{},
	}
}

func (s sessions) Create(code string, name string) string {
	uuid := s.uuidCreator.Create()

	s.mut.Lock()
	s.sessionById[uuid] = playerSession{
		partyCode: code,
		name:      name,
	}
	s.mut.Unlock()

	return uuid
}
