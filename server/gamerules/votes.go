package gamerules

import (
	"errors"
)

var (
	errPlayerHasAlreadyVoted = errors.New("player has already voted")
)

type votes map[string]bool

func (v votes) voteBy(name string, approved bool) (votes, error) {
	if v == nil {
		v = make(votes)
	}

	if a, exists := v[name]; exists && a == approved {
		return v, errPlayerHasAlreadyVoted
	}

	v[name] = approved
	return v, nil
}

func (v votes) approveBy(name string) (votes, error) {
	return v.voteBy(name, true)
}

func (v votes) rejectBy(name string) (votes, error) {
	return v.voteBy(name, false)
}
