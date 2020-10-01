package sessions

type playerSession struct {
	partyCode string
	name      string
}

type uuidCreator interface {
	Create() string
}

type sessions struct {
	uuidCreator
	sessionById map[string]playerSession
}

func New(uuidCreator uuidCreator) sessions {
	return sessions{
		uuidCreator: uuidCreator,
		sessionById: make(map[string]playerSession),
	}
}

func (s sessions) Create(code string, name string) string {
	uuid := s.uuidCreator.Create()

	s.sessionById[uuid] = playerSession{
		partyCode: code,
		name:      name,
	}

	return uuid
}
