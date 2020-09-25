package gamerules

import "errors"

const (
	minNumberOfPlayers = 5
	maxNumberOfPlayers = 10
)

var (
	errPlayerNotFound            = errors.New("player not found")
	errNotEnoughPlayers          = errors.New("need 5 players before starting game")
	errPlayerAlreadyInGame       = errors.New("player already in game")
	errAlreadyMaxNumberOfPlayers = errors.New("can't have more than 10 players")
)

type players []string

func (p players) add(name string) (players, error) {
	if len(p) == maxNumberOfPlayers {
		return p, errAlreadyMaxNumberOfPlayers
	}

	if p.exists(name) {
		return p, errPlayerAlreadyInGame
	}

	return append(p, name), nil
}

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

func (p players) hasMinNumberOfPlayers() bool {
	return len(p) >= minNumberOfPlayers
}
