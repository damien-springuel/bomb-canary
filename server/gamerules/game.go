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

type State string

const (
	NotStarted        State = "notStarted"
	SelectingTeam     State = "selectingTeam"
	VotingOnTeam      State = "votingOnTeam"
	ConductingMission State = "conductingMission"
	GameOver          State = "gameOver"
)

type Mission int

const (
	First  Mission = 1
	Second Mission = 2
	Third  Mission = 3
	Fourth Mission = 4
	Fifth  Mission = 5
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
	missionRequirementsByNumberOfPlayer = map[int]map[Mission]missionRequirement{
		5: {
			First:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			Second: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Third:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			Fourth: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Fifth:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
		},
		6: {
			First:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			Second: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			Fourth: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Fifth:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
		},
		7: {
			First:  {nbOfPeopleToGo: 2, nbFailuresRequiredToFailMission: 1},
			Second: {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Third:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Fourth: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 2},
			Fifth:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
		},
		8: {
			First:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Second: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			Third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			Fourth: {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 2},
			Fifth:  {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 1},
		},
		9: {
			First:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Second: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			Third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			Fourth: {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 2},
			Fifth:  {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 1},
		},
		10: {
			First:  {nbOfPeopleToGo: 3, nbFailuresRequiredToFailMission: 1},
			Second: {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			Third:  {nbOfPeopleToGo: 4, nbFailuresRequiredToFailMission: 1},
			Fourth: {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 2},
			Fifth:  {nbOfPeopleToGo: 5, nbFailuresRequiredToFailMission: 1},
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
	state            State
	players          players
	playerAllegiance map[string]Allegiance
	leader           string
	currentTeam      players
	currentMission   Mission
	teamVotes        votes
	voteFailures     int
	missionOutcomes  votes
	missionResults   missionResults
}

func NewGame() Game {
	return Game{state: NotStarted}
}

func (g Game) AddPlayer(name string) (Game, error) {
	if g.state != NotStarted {
		return g, fmt.Errorf("%w: can only add player during %s state, state was %s", errInvalidStateForAction, NotStarted, g.state)
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
	if g.state != NotStarted {
		return g, fmt.Errorf("%w: can only remove player during %s state, state was %s", errInvalidStateForAction, NotStarted, g.state)
	}

	p, err := g.players.remove(name)
	if err != nil {
		return g, err
	}

	g.players = p
	return g, nil
}

func (g Game) Start(allegianceGenerator AllegianceGenerator) (Game, error) {
	if g.state != NotStarted {
		return g, fmt.Errorf("%w: can only start the game during %s state, state was %s", errInvalidStateForAction, NotStarted, g.state)
	}

	if g.players.count() < minNumberOfPlayers {
		return g, errNotEnoughPlayers
	}

	allegiances := allegianceGenerator.Generate(g.players.count(), nbOfSpiesByNumberOfPlayers[g.players.count()])
	g.playerAllegiance = map[string]Allegiance{}
	for i, a := range allegiances {
		g.playerAllegiance[g.players[i]] = a
	}

	g.state = SelectingTeam
	g.leader = g.players[0]
	g.currentMission = First
	return g, nil
}

func (g Game) nbPeopleThatHaveToGoOnMission() int {
	return missionRequirementsByNumberOfPlayer[g.players.count()][g.currentMission].nbOfPeopleToGo
}

func (g Game) nbFailuresRequiredToFailMission() int {
	return missionRequirementsByNumberOfPlayer[g.players.count()][g.currentMission].nbFailuresRequiredToFailMission
}

func (g Game) LeaderSelectsMember(name string) (Game, error) {
	if g.state != SelectingTeam {
		return g, fmt.Errorf("%w: can only select team members during %s state, state was %s", errInvalidStateForAction, SelectingTeam, g.state)
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

func (g Game) LeaderDeselectsMember(name string) (Game, error) {
	if g.state != SelectingTeam {
		return g, fmt.Errorf("%w: can only deselect team members during %s state, state was %s", errInvalidStateForAction, SelectingTeam, g.state)
	}

	newTeam, err := g.currentTeam.remove(name)
	if err != nil {
		return g, err
	}

	g.currentTeam = newTeam
	return g, nil
}

func (g Game) LeaderConfirmsTeamSelection() (Game, error) {
	if g.state != SelectingTeam {
		return g, fmt.Errorf("%w: can only deselect team members during %s state, state was %s", errInvalidStateForAction, SelectingTeam, g.state)
	}

	if g.currentTeam.count() < g.nbPeopleThatHaveToGoOnMission() {
		return g, fmt.Errorf("%w: need %d people, currently have %d", errTeamIsIncomplete, g.nbPeopleThatHaveToGoOnMission(), g.currentTeam.count())
	}

	g.state = VotingOnTeam
	return g, nil
}

func (g Game) voteBy(name string, voter func(name string) (votes, error)) (Game, error) {
	if g.state != VotingOnTeam {
		return g, fmt.Errorf("%w: can only vote on team during %s state, state was %s", errInvalidStateForAction, VotingOnTeam, g.state)
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
			g.state = ConductingMission
			g.voteFailures = 0
		} else {
			g.state = SelectingTeam
			g.voteFailures += 1

			if g.voteFailures == maxVoteFailures {
				g.state = GameOver
			} else {
				g.leader = g.players.after(g.leader)
			}
		}
		newVotes = nil
	}

	g.teamVotes = newVotes
	return g, nil
}

func (g Game) ApproveTeamBy(name string) (Game, error) {
	return g.voteBy(name, g.teamVotes.approveBy)
}

func (g Game) RejectTeamBy(name string) (Game, error) {
	return g.voteBy(name, g.teamVotes.rejectBy)
}

func (g Game) workOnMissionBy(name string, worker func(name string) (votes, error)) (Game, error) {
	if g.state != ConductingMission {
		return g, fmt.Errorf("%w: can only work on mission during %s state, state was %s", errInvalidStateForAction, ConductingMission, g.state)
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
			g.state = GameOver
		} else {
			g.state = SelectingTeam
			g.currentMission += 1
			g.leader = g.players.after(g.leader)
			g.currentTeam = nil
		}
		newOutcomes = nil
	}

	g.missionOutcomes = newOutcomes
	return g, nil
}

func (g Game) SucceedMissionBy(name string) (Game, error) {
	return g.workOnMissionBy(name, g.missionOutcomes.approveBy)
}

func (g Game) FailMissionBy(name string) (Game, error) {
	return g.workOnMissionBy(name, g.missionOutcomes.rejectBy)
}

func (g Game) Leader() string {
	return g.leader
}

func (g Game) State() State {
	return g.state
}

func (g Game) VoteFailures() int {
	return g.voteFailures
}

func (g Game) GetMissionResults() map[Mission]bool {
	return g.missionResults.copy()
}

func (g Game) CurrentMission() Mission {
	return g.currentMission
}

func (g Game) Winner() Allegiance {
	if g.state != GameOver {
		return ""
	}

	successes := 0
	for _, success := range g.missionResults {
		if success {
			successes += 1
		}
	}

	if successes == 3 {
		return Resistance
	}
	return Spy
}
