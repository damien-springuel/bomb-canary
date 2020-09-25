package gamerules

import (
	"errors"
	"fmt"
)

var (
	errPlayerNotFound        = errors.New("player not found")
	errNotEnoughPlayers      = errors.New("not enough players")
	errInvalidStateForAction = errors.New("invalid state for action")
)

type state string

const (
	notStarted    state = "notStarted"
	selectingTeam state = "selectingTeam"
)

type players []string

func (p players) remove(name string) (players, error) {
	index := 0
	found := false
	for i, n := range p {
		if n == name {
			index = i
			found = true
			break
		}
	}

	if !found {
		return p, errPlayerNotFound
	}

	return append(p[:index], p[index+1:]...), nil
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

	if len(g.players) < 5 {
		return g, errNotEnoughPlayers
	}

	g.state = selectingTeam
	return g, nil
}
