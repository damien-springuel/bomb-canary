package gamerules

import (
	"testing"

	. "github.com/onsi/gomega"
)

func createNewlyStartedGame() game {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")
	newGame, _ = newGame.addPlayer("Edith")
	newGame, _ = newGame.start()
	return newGame
}

func createNewlyVotingOnTeamGame() game {
	newGame := createNewlyStartedGame()
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	return newGame
}

func createNewlyConductingMissionGame() game {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, _ = newGame.approveTeamBy("Edith")
	return newGame
}

func Test_CreateGame(t *testing.T) {
	newGame := newGame()
	g := NewWithT(t)
	g.Expect(newGame).To(Equal(game{state: notStarted}))
}

func Test_AddPlayer(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")

	g := NewWithT(t)
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}

func Test_AddPlayer_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")
	newGame, _ = newGame.addPlayer("Edith")

	newGame, _ = newGame.start()

	newGame, err := newGame.addPlayer("Frank")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_AddPlayer_ShouldErrorIfPlayerAlreadyThere(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, err := newGame.addPlayer("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerAlreadyInGroup))
}

func Test_AddPlayer_ShouldErrorIfAlready10Players(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("1")
	newGame, _ = newGame.addPlayer("2")
	newGame, _ = newGame.addPlayer("3")
	newGame, _ = newGame.addPlayer("4")
	newGame, _ = newGame.addPlayer("5")
	newGame, _ = newGame.addPlayer("6")
	newGame, _ = newGame.addPlayer("7")
	newGame, _ = newGame.addPlayer("8")
	newGame, _ = newGame.addPlayer("9")
	newGame, _ = newGame.addPlayer("10")
	newGame, err := newGame.addPlayer("11")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errAlreadyMaxNumberOfPlayers))
}

func Test_RemovePlayer(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, err := newGame.removePlayer("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.players).To(Equal(players([]string{})))
}

func Test_RemovePlayer_ShouldErrorIfPlayerNotFound(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, err := newGame.removePlayer("Bob")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}

func Test_RemovePlayer_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")
	newGame, _ = newGame.addPlayer("Edith")

	newGame, _ = newGame.start()

	newGame, err := newGame.removePlayer("Bob")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_StartGame_WhenFewerThan5Players_ShouldError(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")

	newGame, err := newGame.start()

	g := NewWithT(t)
	g.Expect(err).To(Equal(errNotEnoughPlayers))
}

func Test_StartGame(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")
	newGame, _ = newGame.addPlayer("Edith")

	newGame, err := newGame.start()

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(selectingTeam))
	g.Expect(newGame.leader).To(Equal("Alice"))
	g.Expect(newGame.currentTeam).To(BeEmpty())
	g.Expect(newGame.currentMission).To(Equal(first))
}

func Test_StartGame_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")
	newGame, _ = newGame.addPlayer("Edith")

	newGame, _ = newGame.start()
	newGame, err := newGame.start()

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_LeaderSelectsAMember(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, err := newGame.leaderSelectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.currentTeam).To(Equal(players([]string{"Alice"})))
}

func Test_LeaderSelectsAMember_ShouldErrorIfPlayerNotExists(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, err := newGame.leaderSelectsMember("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_LeaderSelectsAMember_ShouldErrorIfAlreadyInTeam(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, err := newGame.leaderSelectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerAlreadyInGroup))
}

func Test_LeaderSelectsAMember_ShouldErrorIfTeamIsComplete(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, err := newGame.leaderSelectsMember("Charlie")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errTeamIsFull))
}

func Test_LeaderSelectsAMember_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := newGame()

	newGame, err := newGame.leaderSelectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_LeaderDeselectsAMember(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, err := newGame.leaderDeselectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.currentTeam).To(BeEmpty())
}

func Test_LeaderDeselectsAMember_ShouldErrorIfPlayerNotInTeam(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, err := newGame.leaderDeselectsMember("Bob")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_LeaderDeselectsAMember_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := newGame()

	newGame, err := newGame.leaderDeselectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_LeaderConfirmsSelection(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, err := newGame.leaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(votingOnTeam))
}

func Test_LeaderConfirmsSelection_ShouldErrorIfTeamIsIncomplete(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, err := newGame.leaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errTeamIsIncomplete))
}

func Test_LeaderConfirmsSelection_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := newGame()

	newGame, err := newGame.leaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_ApproveTeam(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.approveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.teamVotes).To(Equal(votes(map[string]bool{"Alice": true})))
	g.Expect(newGame.state).To(Equal(votingOnTeam))
}

func Test_ApproveTeam_ShouldReturnErrorIfAlreadyApproved(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, err := newGame.approveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_ApproveTeam_ShouldReturnErrorIfAlreadyVoted(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, err := newGame.rejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_ApproveTeam_ShouldReturnErrorIfPlayerDoesntExist(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.approveTeamBy("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_ApproveTeam_ShouldReturnErrorIfNotCurrentlyVoting(t *testing.T) {
	newGame := createNewlyStartedGame()
	newGame, err := newGame.approveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_RejectTeam(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.rejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.teamVotes).To(Equal(votes(map[string]bool{"Alice": false})))
	g.Expect(newGame.state).To(Equal(votingOnTeam))
}

func Test_RejectTeam_ShouldReturnErrorIfAlreadyRejected(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.rejectTeamBy("Alice")
	newGame, err := newGame.rejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_RejectTeam_ShouldReturnErrorIfAlreadyVoted(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.rejectTeamBy("Alice")
	newGame, err := newGame.approveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
}

func Test_RejectTeam_ShouldReturnErrorIfPlayerDoesntExist(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, err := newGame.rejectTeamBy("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_RejectTeam_ShouldReturnErrorIfNotCurrentlyVoting(t *testing.T) {
	newGame := createNewlyStartedGame()
	newGame, err := newGame.rejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_ApproveRejectTeam_ShouldMoveToConductingMissionIfVoteHasMajority(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, err := newGame.rejectTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(conductingMission))
}

func Test_ApproveRejectTeam_ShouldMoveToSelectingTeamIfVoteDoesntHaveMajority_AndSwitchLeader(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.rejectTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, err := newGame.rejectTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(selectingTeam))
	g.Expect(newGame.voteFailures).To(Equal(1))
	g.Expect(newGame.leader).To(Equal("Bob"))
}

func Test_ApproveRejectTeam_ShouldMoveToGameOverIfVoteFailed5TimesInARow(t *testing.T) {
	g := NewWithT(t)
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.rejectTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, _ = newGame.rejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Bob"))
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.rejectTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, _ = newGame.rejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Charlie"))
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.rejectTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, _ = newGame.rejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Dan"))
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.rejectTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, _ = newGame.rejectTeamBy("Edith")

	g.Expect(newGame.leader).To(Equal("Edith"))
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.rejectTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, err := newGame.rejectTeamBy("Edith")

	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(gameOver))
	g.Expect(newGame.voteFailures).To(Equal(5))
	g.Expect(newGame.leader).To(Equal("Edith"))
}

func Test_ApproveRejectTeam_VoteFailureShouldResetAfterASuccessfulVote(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.rejectTeamBy("Charlie")
	newGame, _ = newGame.rejectTeamBy("Dan")
	newGame, _ = newGame.rejectTeamBy("Edith")

	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, err := newGame.approveTeamBy("Edith")

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
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderSelectsMember("Charlie")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, _ = newGame.approveTeamBy("Edith")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.succeedMissionBy("Charlie")

	// Third turn
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, _ = newGame.approveTeamBy("Edith")
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
	newGame := newGame()
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")
	newGame, _ = newGame.addPlayer("Edith")
	newGame, _ = newGame.addPlayer("Fred")
	newGame, _ = newGame.addPlayer("Gordon")
	newGame, _ = newGame.start()

	// First turn
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, _ = newGame.approveTeamBy("Edith")
	newGame, _ = newGame.approveTeamBy("Fred")
	newGame, _ = newGame.approveTeamBy("Gordon")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")

	// Second turn
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderSelectsMember("Charlie")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, _ = newGame.approveTeamBy("Edith")
	newGame, _ = newGame.approveTeamBy("Fred")
	newGame, _ = newGame.approveTeamBy("Gordon")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.failMissionBy("Charlie")

	// Third turn
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderSelectsMember("Charlie")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, _ = newGame.approveTeamBy("Edith")
	newGame, _ = newGame.approveTeamBy("Fred")
	newGame, _ = newGame.approveTeamBy("Gordon")
	newGame, _ = newGame.succeedMissionBy("Alice")
	newGame, _ = newGame.succeedMissionBy("Bob")
	newGame, _ = newGame.failMissionBy("Charlie")

	// Fourth turn
	newGame, _ = newGame.leaderSelectsMember("Alice")
	newGame, _ = newGame.leaderSelectsMember("Bob")
	newGame, _ = newGame.leaderSelectsMember("Charlie")
	newGame, _ = newGame.leaderSelectsMember("Dan")
	newGame, _ = newGame.leaderConfirmsTeamSelection()
	newGame, _ = newGame.approveTeamBy("Alice")
	newGame, _ = newGame.approveTeamBy("Bob")
	newGame, _ = newGame.approveTeamBy("Charlie")
	newGame, _ = newGame.approveTeamBy("Dan")
	newGame, _ = newGame.approveTeamBy("Edith")
	newGame, _ = newGame.approveTeamBy("Fred")
	conductingFourthMission, _ := newGame.approveTeamBy("Gordon")
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
