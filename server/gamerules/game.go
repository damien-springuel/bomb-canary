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

type MissionRequirement struct {
	NbOfPeopleToGo                  int
	NbFailuresRequiredToFailMission int
}

var (
	missionRequirementsByNumberOfPlayer = map[int]map[Mission]MissionRequirement{
		5: {
			First:  {NbOfPeopleToGo: 2, NbFailuresRequiredToFailMission: 1},
			Second: {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Third:  {NbOfPeopleToGo: 2, NbFailuresRequiredToFailMission: 1},
			Fourth: {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Fifth:  {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
		},
		6: {
			First:  {NbOfPeopleToGo: 2, NbFailuresRequiredToFailMission: 1},
			Second: {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Third:  {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
			Fourth: {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Fifth:  {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
		},
		7: {
			First:  {NbOfPeopleToGo: 2, NbFailuresRequiredToFailMission: 1},
			Second: {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Third:  {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Fourth: {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 2},
			Fifth:  {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
		},
		8: {
			First:  {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Second: {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
			Third:  {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
			Fourth: {NbOfPeopleToGo: 5, NbFailuresRequiredToFailMission: 2},
			Fifth:  {NbOfPeopleToGo: 5, NbFailuresRequiredToFailMission: 1},
		},
		9: {
			First:  {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Second: {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
			Third:  {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
			Fourth: {NbOfPeopleToGo: 5, NbFailuresRequiredToFailMission: 2},
			Fifth:  {NbOfPeopleToGo: 5, NbFailuresRequiredToFailMission: 1},
		},
		10: {
			First:  {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
			Second: {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
			Third:  {NbOfPeopleToGo: 4, NbFailuresRequiredToFailMission: 1},
			Fourth: {NbOfPeopleToGo: 5, NbFailuresRequiredToFailMission: 2},
			Fifth:  {NbOfPeopleToGo: 5, NbFailuresRequiredToFailMission: 1},
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
	state           State
	players         players
	leader          string
	currentTeam     players
	currentMission  Mission
	teamVotes       votes
	voteFailures    int
	missionOutcomes votes
	missionResults  missionResults
	spies           players
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

func (g Game) Start(allegianceGenerator AllegianceGenerator) (Game, map[string]Allegiance, map[Mission]MissionRequirement, error) {
	if g.state != NotStarted {
		return g, nil, nil, fmt.Errorf("%w: can only start the game during %s state, state was %s", errInvalidStateForAction, NotStarted, g.state)
	}

	if g.players.count() < minNumberOfPlayers {
		return g, nil, nil, errNotEnoughPlayers
	}

	g.state = SelectingTeam
	g.leader = g.players[0]
	g.currentMission = First

	allegiances := allegianceGenerator.Generate(g.players.count(), nbOfSpiesByNumberOfPlayers[g.players.count()])
	playerAllegiance := map[string]Allegiance{}
	for i, a := range allegiances {
		playerAllegiance[g.players[i]] = a
		if a == Spy {
			g.spies, _ = g.spies.add(g.players[i])
		}
	}

	return g, playerAllegiance, g.getMissionRequirements(), nil
}

func (g Game) getMissionRequirements() map[Mission]MissionRequirement {
	requirements := make(map[Mission]MissionRequirement, len(missionRequirementsByNumberOfPlayer[g.players.count()]))
	for mission, req := range missionRequirementsByNumberOfPlayer[g.players.count()] {
		requirements[mission] = req
	}
	return requirements
}

func (g Game) nbPeopleThatHaveToGoOnMission() int {
	return missionRequirementsByNumberOfPlayer[g.players.count()][g.currentMission].NbOfPeopleToGo
}

func (g Game) nbFailuresRequiredToFailMission() int {
	return missionRequirementsByNumberOfPlayer[g.players.count()][g.currentMission].NbFailuresRequiredToFailMission
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

func (g Game) voteBy(name string, voter func(name string) (votes, error)) (Game, map[string]bool, error) {
	if g.state != VotingOnTeam {
		return g, nil, fmt.Errorf("%w: can only vote on team during %s state, state was %s", errInvalidStateForAction, VotingOnTeam, g.state)
	}

	if !g.players.exists(name) {
		return g, nil, errPlayerNotFound
	}

	newVotes, err := voter(name)
	if err != nil {
		return g, nil, err
	}

	if newVotes.hasEveryoneVoted(g.players.count()) {
		if newVotes.hasMajority() {
			g.state = ConductingMission
			g.voteFailures = 0
		} else {
			g.state = SelectingTeam
			g.voteFailures += 1
			g.currentTeam = nil

			if g.voteFailures == maxVoteFailures {
				g.state = GameOver
			} else {
				g.leader = g.players.after(g.leader)
			}
		}
		g.teamVotes = nil
	} else {
		g.teamVotes = newVotes
	}

	return g, newVotes.copy(), nil
}

func (g Game) ApproveTeamBy(name string) (Game, map[string]bool, error) {
	return g.voteBy(name, g.teamVotes.approveBy)
}

func (g Game) RejectTeamBy(name string) (Game, map[string]bool, error) {
	return g.voteBy(name, g.teamVotes.rejectBy)
}

func (g Game) workOnMissionBy(name string, worker func(name string) (votes, error)) (Game, map[string]bool, error) {
	if g.state != ConductingMission {
		return g, nil, fmt.Errorf("%w: can only work on mission during %s state, state was %s", errInvalidStateForAction, ConductingMission, g.state)
	}

	if !g.currentTeam.exists(name) {
		return g, nil, errPlayerNotFound
	}

	newOutcomes, err := worker(name)
	if err != nil {
		return g, nil, err
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
		g.missionOutcomes = nil
	} else {
		g.missionOutcomes = newOutcomes
	}

	return g, newOutcomes.copy(), nil
}

func (g Game) SucceedMissionBy(name string) (Game, map[string]bool, error) {
	return g.workOnMissionBy(name, g.missionOutcomes.approveBy)
}

func (g Game) FailMissionBy(name string) (Game, map[string]bool, error) {
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

	if g.voteFailures > 0 {
		return Spy
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

func (g Game) Spies() []string {
	return g.spies
}
