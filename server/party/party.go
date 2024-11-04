package party

import (
	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type dispatcher interface {
	Dispatch(m messagebus.Message)
}

type partyService struct {
	dispatcher dispatcher
}

func NewPartyService(dispatcher dispatcher) partyService {
	return partyService{
		dispatcher: dispatcher,
	}
}

func (p partyService) JoinParty(name string) {
	p.dispatcher.Dispatch(messagebus.JoinParty{Player: name})
}
