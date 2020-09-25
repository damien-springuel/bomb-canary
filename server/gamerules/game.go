package gamerules

import "errors"

var (
	errPlayerNotFound = errors.New("player not found")
)

type state string

const (
	notStarted state = "notStarted"
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

func (g game) addPlayer(name string) game {
	g.players = append(g.players, name)
	return g
}

func (g game) removePlayer(name string) (game, error) {
	p, err := g.players.remove(name)
	if err != nil {
		return g, err
	}
	g.players = p
	return g, nil
}
