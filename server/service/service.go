package service

import (
	"github.com/damien-springuel/bomb-canary/server/gamerules"
	"github.com/damien-springuel/bomb-canary/server/messagebus"
)

type handler func(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messageToDispatch messagebus.Message)

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
	var handler handler
	switch m.(type) {
	case joinParty:
		handler = s.handleJoinPartyCommand
	case startGame:
		handler = s.handleStartGameCommand
	}
	updatedGame, messageToDispatch := handler(s.gamesByPartyCode[m.GetPartyCode()], m)
	s.gamesByPartyCode[m.GetPartyCode()] = updatedGame
	s.messageDispatcher.dispatchMessage(messageToDispatch)
}

func (s service) getGameForPartyCode(code string) gamerules.Game {
	return s.gamesByPartyCode[code]
}

func (s service) handleJoinPartyCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messageToDispatch messagebus.Message) {
	joinPartyCommand := message.(joinParty)
	updatedGame, _ = currentGame.AddPlayer(joinPartyCommand.user)
	messageToDispatch = playerJoined{party: party{code: message.GetPartyCode()}, user: joinPartyCommand.user}
	return
}

func (s service) handleStartGameCommand(currentGame gamerules.Game, message messagebus.Message) (updatedGame gamerules.Game, messageToDispatch messagebus.Message) {
	updatedGame, _ = currentGame.Start(s.allegianceGenerator)
	messageToDispatch = gameStarted{party: party{code: message.GetPartyCode()}}
	return
}
