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

	if _, exists := v[name]; exists {
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

func (v votes) hasEveryoneVoted(nbPlayers int) bool {
	return len(v) == nbPlayers
}

func (v votes) hasMajority() bool {
	nbVotes := len(v)
	majority := (nbVotes / 2) + 1

	nbApproved := 0
	for _, approved := range v {
		if approved {
			nbApproved += 1
		}
	}

	return nbApproved >= majority
}

func (v votes) nbRejections() int {
	nbRejections := 0

	for _, approved := range v {
		if !approved {
			nbRejections += 1
		}
	}

	return nbRejections
}

func (v votes) copy() votes {
	copy := make(votes)
	for name, vote := range v {
		copy[name] = vote
	}
	return copy
}
