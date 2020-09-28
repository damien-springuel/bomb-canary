package service

import (
	"fmt"

	"github.com/damien-springuel/bomb-canary/server/gamerules"
)

type codeGenerator interface {
	generateCode() string
}

type message interface{}

type messageDispatcher interface {
	dispatchMessage(m message)
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

func (s service) handleMessage(m message) error {
	switch v := m.(type) {
	case joinParty:
		updatedGame, err := s.games[v.partyCode].AddPlayer(v.user)
		if err != nil {
			return fmt.Errorf("can't join party: %w", err)
		}
		s.games[v.partyCode] = updatedGame
		s.messageDispatcher.dispatchMessage(playerJoined{partyCode: v.partyCode, user: v.user})
	}

	return nil
}

func (s service) getGameForPartyCode(code string) gamerules.Game {
	return s.games[code]
}
