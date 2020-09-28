package gamerules

import "errors"

var (
	errPlayerNotFound       = errors.New("player not found")
	errPlayerAlreadyInGroup = errors.New("player already in group")
)

type players []string

func (p players) add(name string) (players, error) {
	if p.exists(name) {
		return p, errPlayerAlreadyInGroup
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

func (p players) count() int {
	return len(p)
}

func (p players) after(name string) string {
	i, _ := p.index(name)
	next := (i + 1) % len(p)
	return p[next]
}
