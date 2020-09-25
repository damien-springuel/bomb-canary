package gamerules

import (
	"errors"
	"fmt"
)

var (
	errInvalidStateForAction = errors.New("invalid state for action")
)

type state string

const (
	notStarted    state = "notStarted"
	selectingTeam state = "selectingTeam"
)

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

	p, err := g.players.add(name)
	if err != nil {
		return g, err
	}

	g.players = p
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

	if !g.players.hasMinNumberOfPlayers() {
		return g, errNotEnoughPlayers
	}

	g.state = selectingTeam
	return g, nil
}
