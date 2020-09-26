package gamerules

import (
	"errors"
)

var (
	errPlayerHasAlreadyVoted = errors.New("player has already voted")
)

type votes map[string]bool

func (v votes) approveBy(name string) (votes, error) {
	if v == nil {
		v = make(votes)
	}

	if approved, exists := v[name]; exists && approved {
		return v, errPlayerHasAlreadyVoted
	}

	v[name] = true

	return v, nil
}
