package gamerules

import (
	"errors"
	"fmt"
)

const (
	minNumberOfPlayers = 5
	maxNumberOfPlayers = 10
)

var (
	errPlayerNotFound            = errors.New("player not found")
	errNotEnoughPlayers          = errors.New("need 5 players before starting game")
	errInvalidStateForAction     = errors.New("invalid state for action")
	errPlayerAlreadyInGame       = errors.New("player already in game")
	errAlreadyMaxNumberOfPlayers = errors.New("can't have more than 10 players")
)

type state string

const (
	notStarted    state = "notStarted"
	selectingTeam state = "selectingTeam"
)

type players []string

func (p players) remove(name string) (players, error) {
	index, exists := p.index(name)

	if !exists {
		return p, errPlayerNotFound
	}

	return append(p[:index], p[index+1:]...), nil
}

func (p players) index(name string) (int, bool) {
	for i, n := range p {
		if n == name {
			return i, true
		}
	}
	return -1, false
}

func (p players) exists(name string) bool {
	_, exists := p.index(name)
	return exists
}

type game struct {
	id      string
	state   state
	players players
}

func newGame(id string) game {
	return game{id: id, state: notStarted}
}

func (g game) addPlayer(name string) (game, error) {
	if g.state != notStarted {
		return g, fmt.Errorf("%w: can only add player during %s state, state was %s", errInvalidStateForAction, notStarted, g.state)
	}

	if len(g.players) == maxNumberOfPlayers {
		return g, errAlreadyMaxNumberOfPlayers
	}

	if g.players.exists(name) {
		return g, errPlayerAlreadyInGame
	}

	g.players = append(g.players, name)
	return g, nil
}

func (g game) removePlayer(name string) (game, error) {
	if g.state != notStarted {
		return g, fmt.Errorf("%w: can only remove player during %s state, state was %s", errInvalidStateForAction, notStarted, g.state)
	}

	p, err := g.players.remove(name)
	if err != nil {
		return g, err
	}

	g.players = p
	return g, nil
}

func (g game) start() (game, error) {
	if g.state != notStarted {
		return g, fmt.Errorf("%w: can only start the game during %s state, state was %s", errInvalidStateForAction, notStarted, g.state)
	}

	if len(g.players) < minNumberOfPlayers {
		return g, errNotEnoughPlayers
	}

	g.state = selectingTeam
	return g, nil
}
