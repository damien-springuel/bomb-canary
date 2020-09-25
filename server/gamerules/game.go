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
	errInvalidStateForAction     = errors.New("invalid state for action")
	errAlreadyMaxNumberOfPlayers = fmt.Errorf("can't have more than %d players in game", maxNumberOfPlayers)
	errNotEnoughPlayers          = fmt.Errorf("new at least %d players in game", minNumberOfPlayers)
	errTeamIsFull                = errors.New("team is maxed out")
)

type state string

const (
	notStarted    state = "notStarted"
	selectingTeam state = "selectingTeam"
)

type mission int

const (
	first  mission = 1
	second mission = 2
	third  mission = 3
	fourth mission = 4
	fifth  mission = 5
)

type missionRequirement struct {
	numberOfPeopleToGo int
}

var (
	missionRequirementsByNumberOfPlayer = map[int]map[mission]missionRequirement{
		5: {
			first:  {numberOfPeopleToGo: 2},
			second: {numberOfPeopleToGo: 3},
			third:  {numberOfPeopleToGo: 2},
			fourth: {numberOfPeopleToGo: 3},
			fifth:  {numberOfPeopleToGo: 3},
		},
	}
)

type game struct {
	state          state
	players        players
	leader         string
	currentTeam    players
	currentMission mission
}

func newGame() game {
	return game{state: notStarted}
}

func (g game) addPlayer(name string) (game, error) {
	if g.state != notStarted {
		return g, fmt.Errorf("%w: can only add player during %s state, state was %s", errInvalidStateForAction, notStarted, g.state)
	}

	if g.players.count() == maxNumberOfPlayers {
		return g, errAlreadyMaxNumberOfPlayers
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

	if g.players.count() < minNumberOfPlayers {
		return g, errNotEnoughPlayers
	}

	g.state = selectingTeam
	g.leader = g.players[0]
	g.currentMission = first
	return g, nil
}

func (g game) numberPeopleThatHaveToGoOnNextMission() int {
	return missionRequirementsByNumberOfPlayer[g.players.count()][g.currentMission].numberOfPeopleToGo
}

func (g game) leaderSelectsMember(name string) (game, error) {
	if !g.players.exists(name) {
		return g, errPlayerNotFound
	}

	if len(g.currentTeam) == g.numberPeopleThatHaveToGoOnNextMission() {
		return g, fmt.Errorf("%w: can't have more than %d", errTeamIsFull, g.numberPeopleThatHaveToGoOnNextMission())
	}

	newTeam, err := g.currentTeam.add(name)
	if err != nil {
		return g, err
	}

	g.currentTeam = newTeam
	return g, nil
}

func (g game) leaderDeselectsMember(name string) (game, error) {
	newTeam, err := g.currentTeam.remove(name)
	if err != nil {
		return g, err
	}

	g.currentTeam = newTeam

	return g, nil
}
