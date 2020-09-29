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
	codeGenerator       codeGenerator
	messageDispatcher   messageDispatcher
	allegianceGenerator gamerules.AllegianceGenerator
	gamesByPartyCode    map[string]gamerules.Game
}

func newService(codeGenerator codeGenerator, messageDispatcher messageDispatcher, allegianceGenerator gamerules.AllegianceGenerator) service {
	return service{
		codeGenerator:       codeGenerator,
		messageDispatcher:   messageDispatcher,
		allegianceGenerator: allegianceGenerator,
		gamesByPartyCode:    make(map[string]gamerules.Game),
	}
}

func (s service) createParty() string {
	newCode := s.codeGenerator.generateCode()
	s.gamesByPartyCode[newCode] = gamerules.NewGame()
	return newCode
}

func (s service) handleMessage(m messagebus.Message) {
	switch v := m.(type) {
	case joinParty:
		updatedGame, _ := s.gamesByPartyCode[m.GetPartyCode()].AddPlayer(v.user)
		s.gamesByPartyCode[m.GetPartyCode()] = updatedGame
		s.messageDispatcher.dispatchMessage(playerJoined{party: party{code: m.GetPartyCode()}, user: v.user})
	case startGame:
		updatedGame, _ := s.gamesByPartyCode[m.GetPartyCode()].Start(s.allegianceGenerator)
		s.gamesByPartyCode[m.GetPartyCode()] = updatedGame
		s.messageDispatcher.dispatchMessage(gameStarted{party: party{code: m.GetPartyCode()}})
	}
}

func (s service) getGameForPartyCode(code string) gamerules.Game {
	return s.gamesByPartyCode[code]
}
