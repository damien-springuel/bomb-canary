package gamerules

import (
	"testing"

	. "github.com/onsi/gomega"
)

type spiesFirstGenerator struct{}

func (s spiesFirstGenerator) Generate(nbPlayers, nbSpies int) []Allegiance {
	allegiances := make([]Allegiance, nbPlayers)
	for i := range allegiances {
		if i < nbSpies {
			allegiances[i] = Spy
		} else {
			allegiances[i] = Resistance
		}
	}
	return allegiances
}

func createNewlyStartedGame() Game {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")
	newGame, _ = newGame.AddPlayer("Edith")
	newGame, _ = newGame.Start(spiesFirstGenerator{})
	return newGame
}

func createNewlyVotingOnTeamGame() Game {
	newGame := createNewlyStartedGame()
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	return newGame
}

func createNewlyConductingMissionGame() Game {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, _ = newGame.ApproveTeamBy("Edith")
	return newGame
}

func Test_CreateGame(t *testing.T) {
	newGame := NewGame()
	g := NewWithT(t)
	g.Expect(newGame).To(Equal(Game{state: notStarted}))
}

func Test_AddPlayer(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")

	g := NewWithT(t)
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}

func Test_AddPlayer_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")
	newGame, _ = newGame.AddPlayer("Edith")

	newGame, _ = newGame.Start(spiesFirstGenerator{})

	newGame, err := newGame.AddPlayer("Frank")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_AddPlayer_ShouldErrorIfPlayerAlreadyThere(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, err := newGame.AddPlayer("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerAlreadyInGroup))
}

func Test_AddPlayer_ShouldErrorIfAlready10Players(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("1")
	newGame, _ = newGame.AddPlayer("2")
	newGame, _ = newGame.AddPlayer("3")
	newGame, _ = newGame.AddPlayer("4")
	newGame, _ = newGame.AddPlayer("5")
	newGame, _ = newGame.AddPlayer("6")
	newGame, _ = newGame.AddPlayer("7")
	newGame, _ = newGame.AddPlayer("8")
	newGame, _ = newGame.AddPlayer("9")
	newGame, _ = newGame.AddPlayer("10")
	newGame, err := newGame.AddPlayer("11")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errAlreadyMaxNumberOfPlayers))
}

func Test_RemovePlayer(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, err := newGame.removePlayer("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.players).To(Equal(players([]string{})))
}

func Test_RemovePlayer_ShouldErrorIfPlayerNotFound(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, err := newGame.removePlayer("Bob")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}

func Test_RemovePlayer_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")
	newGame, _ = newGame.AddPlayer("Edith")

	newGame, _ = newGame.Start(spiesFirstGenerator{})

	newGame, err := newGame.removePlayer("Bob")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_StartGame_WhenFewerThan5Players_ShouldError(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")

	newGame, err := newGame.Start(spiesFirstGenerator{})

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errNotEnoughPlayers))
}

type spyGenerator struct {
	nbPlayersGiven, nbSpiesGiven int
}

func (s *spyGenerator) Generate(nbPlayers, nbSpies int) []Allegiance {
	s.nbPlayersGiven = nbPlayers
	s.nbSpiesGiven = nbSpies
	return []Allegiance{Spy, Resistance, Spy, Resistance, Resistance}

}

func Test_StartGame(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")
	newGame, _ = newGame.AddPlayer("Edith")

	spyGenerator := &spyGenerator{}
	newGame, err := newGame.Start(spyGenerator)

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(selectingTeam))
	g.Expect(newGame.leader).To(Equal("Alice"))
	g.Expect(newGame.currentTeam).To(BeEmpty())
	g.Expect(newGame.currentMission).To(Equal(first))
	g.Expect(newGame.playerAllegiance).To(Equal(map[string]Allegiance{
		"Alice":   Spy,
		"Bob":     Resistance,
		"Charlie": Spy,
		"Dan":     Resistance,
		"Edith":   Resistance,
	}))
	g.Expect(spyGenerator.nbPlayersGiven).To(Equal(5))
	g.Expect(spyGenerator.nbSpiesGiven).To(Equal(2))
}

func Test_StartGame_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")
	newGame, _ = newGame.AddPlayer("Edith")

	newGame, _ = newGame.Start(spiesFirstGenerator{})
	newGame, err := newGame.Start(spiesFirstGenerator{})

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_LeaderSelectsAMember(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, err := newGame.LeaderSelectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.currentTeam).To(Equal(players([]string{"Alice"})))
}

func Test_LeaderSelectsAMember_ShouldErrorIfPlayerNotExists(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, err := newGame.LeaderSelectsMember("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_LeaderSelectsAMember_ShouldErrorIfAlreadyInTeam(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, err := newGame.LeaderSelectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerAlreadyInGroup))
}

func Test_LeaderSelectsAMember_ShouldErrorIfTeamIsComplete(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, err := newGame.LeaderSelectsMember("Charlie")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errTeamIsFull))
}

func Test_LeaderSelectsAMember_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := NewGame()

	newGame, err := newGame.LeaderSelectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_LeaderDeselectsAMember(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, err := newGame.LeaderDeselectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.currentTeam).To(BeEmpty())
}

func Test_LeaderDeselectsAMember_ShouldErrorIfPlayerNotInTeam(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, err := newGame.LeaderDeselectsMember("Bob")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_LeaderDeselectsAMember_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := NewGame()

	newGame, err := newGame.LeaderDeselectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_LeaderConfirmsSelection(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, err := newGame.LeaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(votingOnTeam))
}

func Test_LeaderConfirmsSelection_ShouldErrorIfTeamIsIncomplete(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, err := newGame.LeaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errTeamIsIncomplete))
}

func Test_LeaderConfirmsSelection_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := NewGame()

	newGame, err := newGame.LeaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_ApproveTeam(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.teamVotes).To(Equal(votes(map[string]bool{"Alice": true})))
	g.Expect(newGame.state).To(Equal(votingOnTeam))
}

func Test_ApproveTeam_ShouldReturnErrorIfAlreadyApproved(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_ApproveTeam_ShouldReturnErrorIfAlreadyVoted(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_ApproveTeam_ShouldReturnErrorIfPlayerDoesntExist(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.ApproveTeamBy("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_ApproveTeam_ShouldReturnErrorIfNotCurrentlyVoting(t *testing.T) {
	newGame := createNewlyStartedGame()
	newGame, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_RejectTeam(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.teamVotes).To(Equal(votes(map[string]bool{"Alice": false})))
	g.Expect(newGame.state).To(Equal(votingOnTeam))
}

func Test_RejectTeam_ShouldReturnErrorIfAlreadyRejected(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.RejectTeamBy("Alice")
	newGame, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_RejectTeam_ShouldReturnErrorIfAlreadyVoted(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.RejectTeamBy("Alice")
	newGame, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_RejectTeam_ShouldReturnErrorIfPlayerDoesntExist(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.RejectTeamBy("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_RejectTeam_ShouldReturnErrorIfNotCurrentlyVoting(t *testing.T) {
	newGame := createNewlyStartedGame()
	newGame, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_ApproveRejectTeam_ShouldMoveToConductingMissionIfVoteHasMajority(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, err := newGame.RejectTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(conductingMission))
}

func Test_ApproveRejectTeam_ShouldMoveToSelectingTeamIfVoteDoesntHaveMajority_AndSwitchLeader(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.RejectTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, err := newGame.RejectTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(selectingTeam))
	g.Expect(newGame.voteFailures).To(Equal(1))
	g.Expect(newGame.leader).To(Equal("Bob"))
}

func Test_ApproveRejectTeam_ShouldMoveToGameOverIfVoteFailed5TimesInARow(t *testing.T) {
	g := NewWithT(t)
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.RejectTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Bob"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.RejectTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Charlie"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.RejectTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Dan"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.RejectTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Edith"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.RejectTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, err := newGame.RejectTeamBy("Edith")

	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(gameOver))
	g.Expect(newGame.voteFailures).To(Equal(5))
	g.Expect(newGame.leader).To(Equal("Edith"))
}

func Test_ApproveRejectTeam_VoteFailureShouldResetAfterASuccessfulVote(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.RejectTeamBy("Charlie")
	newGame, _ = newGame.RejectTeamBy("Dan")
	newGame, _ = newGame.RejectTeamBy("Edith")

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, err := newGame.ApproveTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(conductingMission))
	g.Expect(newGame.voteFailures).To(Equal(0))
}

func Test_SucceedMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, err := newGame.succeedMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.missionOutcomes).To(Equal(votes(map[string]bool{"Alice": true})))
}

func Test_SucceedMission_ShouldFailIfPersonIsntOnTeam(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, err := newGame.succeedMissionBy("Charlie")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
}

func Test_SucceedMission_ShouldFailIfPersonAlreadyWorkedOnMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, err := newGame.succeedMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerHasAlreadyVoted))
}

func Test_SucceedMission_ShouldFailIfStateIsntConductingMission(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.succeedMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_FailMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, err := newGame.failMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.missionOutcomes).To(Equal(votes(map[string]bool{"Alice": false})))
}

func Test_FailMission_ShouldFailIfPersonIsntOnTeam(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, err := newGame.failMissionBy("Charlie")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
}

func Test_FailMission_ShouldFailIfPersonAlreadyWorkedOnMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _ = newGame.failMissionBy("Alice")
	newGame, err := newGame.failMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerHasAlreadyVoted))
}

func Test_FailMission_ShouldFailIfStateIsntConductingMission(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.failMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_SucceedFailMission_ShouldMoveToSelectingTeamWhenEveryoneWorkedOnTheMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _ = newGame.failMissionBy("Alice")
	newGame, err := newGame.failMissionBy("Bob")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(selectingTeam))
	g.Expect(newGame.currentMission).To(Equal(second))
	g.Expect(newGame.missionOutcomes).To(BeNil())
	g.Expect(newGame.missionResults).To(Equal(missionResults(map[mission]bool{first: false})))
	g.Expect(newGame.leader).To(Equal("Bob"))
}

func Test_SucceedFailMission_ShouldMoveToGameOverIfThirdSuccess(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")

	// Second turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, _ = newGame.ApproveTeamBy("Edith")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.succeedMissionBy("Charlie")

	// Third turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, _ = newGame.ApproveTeamBy("Edith")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, err := newGame.succeedMissionBy("Bob")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(gameOver))
	g.Expect(newGame.currentMission).To(Equal(third))
	g.Expect(newGame.missionOutcomes).To(BeNil())
	g.Expect(newGame.missionResults).To(Equal(missionResults(map[mission]bool{first: true, second: true, third: true})))
	g.Expect(newGame.leader).To(Equal("Charlie"))
}

func Test_SucceedFailMission_ShouldSometimesNeedTwoFailureToFailTheMission(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")
	newGame, _ = newGame.AddPlayer("Edith")
	newGame, _ = newGame.AddPlayer("Fred")
	newGame, _ = newGame.AddPlayer("Gordon")
	newGame, _ = newGame.Start(spiesFirstGenerator{})

	// First turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, _ = newGame.ApproveTeamBy("Edith")
	newGame, _ = newGame.ApproveTeamBy("Fred")
	newGame, _ = newGame.ApproveTeamBy("Gordon")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")

	// Second turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, _ = newGame.ApproveTeamBy("Edith")
	newGame, _ = newGame.ApproveTeamBy("Fred")
	newGame, _ = newGame.ApproveTeamBy("Gordon")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.failMissionBy("Charlie")

	// Third turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, _ = newGame.ApproveTeamBy("Edith")
	newGame, _ = newGame.ApproveTeamBy("Fred")
	newGame, _ = newGame.ApproveTeamBy("Gordon")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.failMissionBy("Charlie")

	// Fourth turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderSelectsMember("Dan")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _ = newGame.ApproveTeamBy("Alice")
	newGame, _ = newGame.ApproveTeamBy("Bob")
	newGame, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _ = newGame.ApproveTeamBy("Dan")
	newGame, _ = newGame.ApproveTeamBy("Edith")
	newGame, _ = newGame.ApproveTeamBy("Fred")
	conductingFourthMission, _ := newGame.ApproveTeamBy("Gordon")
	newGame, _ = conductingFourthMission.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.succeedMissionBy("Charlie")
	newGame, _ = newGame.failMissionBy("Dan")

	g := NewWithT(t)
	g.Expect(newGame.missionResults).To(Equal(
		missionResults(map[mission]bool{
			first:  true,
			second: false,
			third:  false,
			fourth: true,
		})))

	newGame, _ = conductingFourthMission.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.failMissionBy("Charlie")
	newGame, _ = newGame.failMissionBy("Dan")
	g.Expect(newGame.missionResults).To(Equal(
		missionResults(map[mission]bool{
			first:  true,
			second: false,
			third:  false,
			fourth: false,
		})))
}
