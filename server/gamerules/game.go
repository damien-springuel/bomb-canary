package gamerules

import (
	"errors"
	"fmt"
)

const (
	minNumberOfPlayers = 5
	maxNumberOfPlayers = 10
	maxVoteFailures    = 5
)

var (
	errInvalidStateForAction     = errors.New("invalid state for action")
	errAlreadyMaxNumberOfPlayers = fmt.Errorf("can't have more than %d players in game", maxNumberOfPlayers)
	errNotEnoughPlayers          = fmt.Errorf("new at least %d players in game", minNumberOfPlayers)
	errTeamIsFull                = errors.New("team is maxed out")
	errTeamIsIncomplete          = errors.New("team is imcomplete")
)

type state string

const (
	notStarted        state = "notStarted"
	selectingTeam     state = "selectingTeam"
	votingOnTeam      state = "votingOnTeam"
	conductingMission state = "conductingMission"
	gameOver          state = "gameOver"
)

type mission int

const (
	first  mission = 1
	second mission = 2
	third  mission = 3
	fourth mission = 4
	fifth  mission = 5
)

type Allegiance string

const (
	Resistance Allegiance = "resistance"
	Spy        Allegiance = "spy"
)

type AllegianceGenerator interface {
	Generate(nbPlayers, nbSpies int) []Allegiance
}

type missionRequirement struct {
	nbOfPeopleToGo                  int
	nbFailuresRequiredToFailMission int
}

var (
	missionRequirementsByNumberOfPlayer = map[int]map[mission]missionRequirement{
		5: {
			first:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			second: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			third:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			fourth: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			fifth:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
		},
		6: {
			first:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			second: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			fourth: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			fifth:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
		},
		7: {
			first:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			second: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			third:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			fourth: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 2},
			fifth:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
		},
		8: {
			first:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			second: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			fourth: {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 2},
			fifth:  {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 1},
		},
		9: {
			first:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			second: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			fourth: {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 2},
			fifth:  {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 1},
		},
		10: {
			first:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			second: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			fourth: {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 2},
			fifth:  {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 1},
		},
	}

	nbOfSpiesByNumberOfPlayers = map[int]int{
		5:  2,
		6:  2,
		7:  3,
		8:  3,
		9:  3,
		10: 4,
	}
)

type Game struct {
	state            state
	players          players
	playerAllegiance map[string]Allegiance
	leader           string
	currentTeam      players
	currentMission   mission
	teamVotes        votes
	voteFailures     int
	missionOutcomes  votes
	missionResults   missionResults
}

func NewGame() Game {
	return Game{state: notStarted}
}

func (g Game) AddPlayer(name string) (Game, error) {
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

func (g Game) removePlayer(name string) (Game, error) {
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

func (g Game) Start(allegianceGenerator AllegianceGenerator) (Game, error) {
	if g.state != notStarted {
		return g, fmt.Errorf("%w: can only start the game during %s state, state was %s", errInvalidStateForAction, notStarted, g.state)
	}

	if g.players.count() < minNumberOfPlayers {
		return g, errNotEnoughPlayers
	}

	allegiances := allegianceGenerator.Generate(g.players.count(), nbOfSpiesByNumberOfPlayers[g.players.count()])
	g.playerAllegiance = map[string]Allegiance{}
	for i, a := range allegiances {
		g.playerAllegiance[g.players[i]] = a
	}

	g.state = selectingTeam
	g.leader = g.players[0]
	g.currentMission = first
	return g, nil
}

func (g Game) nbPeopleThatHaveToGoOnMission() int {
	return missionRequirementsByNumberOfPlayer[g.players.count()][g.currentMission].nbOfPeopleToGo
}

func (g Game) nbFailuresRequiredToFailMission() int {
	return missionRequirementsByNumberOfPlayer[g.players.count()][g.currentMission].nbFailuresRequiredToFailMission
}

func (g Game) leaderSelectsMember(name string) (Game, error) {
	if g.state != selectingTeam {
		return g, fmt.Errorf("%w: can only select team members during %s state, state was %s", errInvalidStateForAction, selectingTeam, g.state)
	}

	if !g.players.exists(name) {
		return g, errPlayerNotFound
	}

	if len(g.currentTeam) == g.nbPeopleThatHaveToGoOnMission() {
		return g, fmt.Errorf("%w: can't have more than %d", errTeamIsFull, g.nbPeopleThatHaveToGoOnMission())
	}

	newTeam, err := g.currentTeam.add(name)
	if err != nil {
		return g, err
	}

	g.currentTeam = newTeam
	return g, nil
}

func (g Game) leaderDeselectsMember(name string) (Game, error) {
	if g.state != selectingTeam {
		return g, fmt.Errorf("%w: can only deselect team members during %s state, state was %s", errInvalidStateForAction, selectingTeam, g.state)
	}

	newTeam, err := g.currentTeam.remove(name)
	if err != nil {
		return g, err
	}

	g.currentTeam = newTeam

	return g, nil
}

func (g Game) leaderConfirmsTeamSelection() (Game, error) {
	if g.state != selectingTeam {
		return g, fmt.Errorf("%w: can only deselect team members during %s state, state was %s", errInvalidStateForAction, selectingTeam, g.state)
	}

	if g.currentTeam.count() < g.nbPeopleThatHaveToGoOnMission() {
		return g, fmt.Errorf("%w: need %d people, currently have %d", errTeamIsIncomplete, g.nbPeopleThatHaveToGoOnMission(), g.currentTeam.count())
	}

	g.state = votingOnTeam
	return g, nil
}

func (g Game) voteBy(name string, voter func(name string) (votes, error)) (Game, error) {
	if g.state != votingOnTeam {
		return g, fmt.Errorf("%w: can only vote on team during %s state, state was %s", errInvalidStateForAction, votingOnTeam, g.state)
	}

	if !g.players.exists(name) {
		return g, errPlayerNotFound
	}

	newVotes, err := voter(name)
	if err != nil {
		return g, err
	}

	if newVotes.hasEveryoneVoted(g.players.count()) {
		if newVotes.hasMajority() {
			g.state = conductingMission
			g.voteFailures = 0
		} else {
			g.state = selectingTeam
			g.voteFailures += 1

			if g.voteFailures == maxVoteFailures {
				g.state = gameOver
			} else {
				g.leader = g.players.after(g.leader)
			}
		}
		newVotes = nil
	}

	g.teamVotes = newVotes
	return g, nil
}

func (g Game) approveTeamBy(name string) (Game, error) {
	return g.voteBy(name, g.teamVotes.approveBy)
}

func (g Game) rejectTeamBy(name string) (Game, error) {
	return g.voteBy(name, g.teamVotes.rejectBy)
}

func (g Game) workOnMissionBy(name string, worker func(name string) (votes, error)) (Game, error) {
	if g.state != conductingMission {
		return g, fmt.Errorf("%w: can only work on mission during %s state, state was %s", errInvalidStateForAction, conductingMission, g.state)
	}

	if !g.currentTeam.exists(name) {
		return g, errPlayerNotFound
	}

	newOutcomes, err := worker(name)
	if err != nil {
		return g, err
	}

	if newOutcomes.hasEveryoneVoted(g.nbPeopleThatHaveToGoOnMission()) {
		if newOutcomes.nbRejections() >= g.nbFailuresRequiredToFailMission() {
			g.missionResults = g.missionResults.failMission(g.currentMission)
		} else {
			g.missionResults = g.missionResults.succeedMission(g.currentMission)
		}

		if g.missionResults.hasThreeSuccessesOrFailures() {
			g.state = gameOver
		} else {
			g.state = selectingTeam
			g.currentMission += 1
			g.leader = g.players.after(g.leader)
		}
		newOutcomes = nil
	}

	g.missionOutcomes = newOutcomes
	return g, nil
}

func (g Game) succeedMissionBy(name string) (Game, error) {
	return g.workOnMissionBy(name, g.missionOutcomes.approveBy)
}

func (g Game) failMissionBy(name string) (Game, error) {
	return g.workOnMissionBy(name, g.missionOutcomes.rejectBy)
}
