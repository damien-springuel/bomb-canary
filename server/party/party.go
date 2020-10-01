package party

import (
	"fmt"

	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type dispatcher interface {
	Dispatch(m messagebus.Message)
}

type partyCreator interface {
	CreateParty() string
}

type partyChecker interface {
	DoesPartyExist(code string) bool
}

type partyService struct {
	creator    partyCreator
	checker    partyChecker
	dispatcher dispatcher
}

func NewPartyService(creator partyCreator, checker partyChecker, dispatcher dispatcher) partyService {
	return partyService{
		creator:    creator,
		checker:    checker,
		dispatcher: dispatcher,
	}
}

func (p partyService) CreateParty() string {
	return p.creator.CreateParty()
}

func (p partyService) JoinParty(code string, name string) error {
	partyExists := p.checker.DoesPartyExist(code)
	if partyExists {
		p.dispatcher.Dispatch(messagebus.JoinParty{Party: messagebus.Party{Code: code}, User: name})
		return nil
	}
	return fmt.Errorf("can't join party {%s}: %w", code, errPartyDoesntExists)
}
