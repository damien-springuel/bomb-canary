package service

import (
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type codeGenerator interface {
	generateCode() string
}

type messageDispatcher interface {
	dispatchMessage(m messagebus.Message)
}

type service struct {
	codeGenerator     codeGenerator
	messageDispatcher messageDispatcher
	games             map[string]gamerules.Game
}

func newService(codeGenerator codeGenerator, messageDispatcher messageDispatcher) service {
	return service{
		codeGenerator:     codeGenerator,
		messageDispatcher: messageDispatcher,
		games:             make(map[string]gamerules.Game),
	}
}

func (s service) createParty() string {
	newCode := s.codeGenerator.generateCode()
	s.games[newCode] = gamerules.NewGame()
	return newCode
}

func (s service) handleMessage(m messagebus.Message) {
	switch v := m.(type) {
	case joinParty:
		updatedGame, err := s.games[v.partyCode].AddPlayer(v.user)
		if err != nil {
			return
		}
		s.games[v.partyCode] = updatedGame
		s.messageDispatcher.dispatchMessage(playerJoined{partyCode: v.partyCode, user: v.user})
	}
}

func (s service) getGameForPartyCode(code string) gamerules.Game {
	return s.games[code]
}
