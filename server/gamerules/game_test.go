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
	newGame, _, _, _ = newGame.Start(spiesFirstGenerator{})
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
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	return newGame
}

func Test_CreateGame(t *testing.T) {
	newGame := NewGame()
	g := NewWithT(t)
	g.Expect(newGame).To(Equal(Game{state: NotStarted}))
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

	newGame, _, _, _ = newGame.Start(spiesFirstGenerator{})

	_, err := newGame.AddPlayer("Frank")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_AddPlayer_ShouldErrorIfPlayerAlreadyThere(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	_, err := newGame.AddPlayer("Alice")

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
	_, err := newGame.AddPlayer("11")

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

	newGame, _, _, _ = newGame.Start(spiesFirstGenerator{})

	_, err := newGame.removePlayer("Bob")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_StartGame_WhenFewerThan5Players_ShouldError(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")

	_, _, _, err := newGame.Start(spiesFirstGenerator{})

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
	newGame, actualPlayerAllegiance, actualMissionRequirements, err := newGame.Start(spyGenerator)

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(SelectingTeam))
	g.Expect(newGame.Leader()).To(Equal("Alice"))
	g.Expect(newGame.currentTeam).To(BeEmpty())
	g.Expect(newGame.CurrentMission()).To(Equal(First))
	g.Expect(actualPlayerAllegiance).To(Equal(map[string]Allegiance{
		"Alice":   Spy,
		"Bob":     Resistance,
		"Charlie": Spy,
		"Dan":     Resistance,
		"Edith":   Resistance,
	}))
	g.Expect(actualMissionRequirements).To(Equal(map[Mission]MissionRequirement{
		First:  {NbOfPeopleToGo: 2, NbFailuresRequiredToFailMission: 1},
		Second: {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
		Third:  {NbOfPeopleToGo: 2, NbFailuresRequiredToFailMission: 1},
		Fourth: {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
		Fifth:  {NbOfPeopleToGo: 3, NbFailuresRequiredToFailMission: 1},
	}))
	g.Expect(spyGenerator.nbPlayersGiven).To(Equal(5))
	g.Expect(spyGenerator.nbSpiesGiven).To(Equal(2))
	g.Expect(newGame.spies).To(Equal(players{"Alice", "Charlie"}))
}

func Test_StartGame_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := NewGame()
	newGame, _ = newGame.AddPlayer("Alice")
	newGame, _ = newGame.AddPlayer("Bob")
	newGame, _ = newGame.AddPlayer("Charlie")
	newGame, _ = newGame.AddPlayer("Dan")
	newGame, _ = newGame.AddPlayer("Edith")

	newGame, _, _, _ = newGame.Start(spiesFirstGenerator{})
	_, _, _, err := newGame.Start(spiesFirstGenerator{})

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

	_, err := newGame.LeaderSelectsMember("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_LeaderSelectsAMember_ShouldErrorIfAlreadyInTeam(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	_, err := newGame.LeaderSelectsMember("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerAlreadyInGroup))
}

func Test_LeaderSelectsAMember_ShouldErrorIfTeamIsComplete(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	_, err := newGame.LeaderSelectsMember("Charlie")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errTeamIsFull))
}

func Test_LeaderSelectsAMember_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := NewGame()

	_, err := newGame.LeaderSelectsMember("Alice")

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
	_, err := newGame.LeaderDeselectsMember("Bob")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
}

func Test_LeaderDeselectsAMember_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := NewGame()

	_, err := newGame.LeaderDeselectsMember("Alice")

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
	g.Expect(newGame.State()).To(Equal(VotingOnTeam))
}

func Test_LeaderConfirmsSelection_ShouldErrorIfTeamIsIncomplete(t *testing.T) {
	newGame := createNewlyStartedGame()

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	_, err := newGame.LeaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errTeamIsIncomplete))
}

func Test_LeaderConfirmsSelection_ShouldErrorIfNotSelectingTeam(t *testing.T) {
	newGame := NewGame()

	_, err := newGame.LeaderConfirmsTeamSelection()

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
}

func Test_ApproveTeam(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, resultingVotes, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.teamVotes).To(Equal(votes(map[string]bool{"Alice": true})))
	g.Expect(resultingVotes).To(Equal(map[string]bool{"Alice": true}))
	g.Expect(newGame.State()).To(Equal(VotingOnTeam))
}

func Test_ApproveTeam_ShouldReturnErrorIfAlreadyApproved(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	_, resultingVotes, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_ApproveTeam_ShouldReturnErrorIfAlreadyVoted(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	_, resultingVotes, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_ApproveTeam_ShouldReturnErrorIfPlayerDoesntExist(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	_, resultingVotes, err := newGame.ApproveTeamBy("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_ApproveTeam_ShouldReturnErrorIfNotCurrentlyVoting(t *testing.T) {
	newGame := createNewlyStartedGame()
	_, resultingVotes, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_RejectTeam(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, resultingVotes, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.teamVotes).To(Equal(votes(map[string]bool{"Alice": false})))
	g.Expect(resultingVotes).To(Equal(map[string]bool{"Alice": false}))
	g.Expect(newGame.State()).To(Equal(VotingOnTeam))
}

func Test_RejectTeam_ShouldReturnErrorIfAlreadyRejected(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.RejectTeamBy("Alice")
	_, resultingVotes, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_RejectTeam_ShouldReturnErrorIfAlreadyVoted(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.RejectTeamBy("Alice")
	_, resultingVotes, err := newGame.ApproveTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerHasAlreadyVoted))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_RejectTeam_ShouldReturnErrorIfPlayerDoesntExist(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	_, resultingVotes, err := newGame.RejectTeamBy("NotThere")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errPlayerNotFound))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_RejectTeam_ShouldReturnErrorIfNotCurrentlyVoting(t *testing.T) {
	newGame := createNewlyStartedGame()
	_, resultingVotes, err := newGame.RejectTeamBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
	g.Expect(resultingVotes).To(BeNil())
}

func Test_ApproveRejectTeam_ShouldMoveToConductingMissionIfVoteHasMajority(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, resultingVotes, err := newGame.RejectTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(ConductingMission))
	g.Expect(resultingVotes).To(Equal(map[string]bool{"Alice": true, "Bob": true, "Charlie": true, "Dan": false, "Edith": false}))
	g.Expect(newGame.teamVotes).To(BeNil())
}

func Test_ApproveRejectTeam_ShouldMoveToSelectingTeamIfVoteDoesntHaveMajority_AndSwitchLeader(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.RejectTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, resultingVotes, err := newGame.RejectTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(SelectingTeam))
	g.Expect(newGame.VoteFailures()).To(Equal(1))
	g.Expect(newGame.Leader()).To(Equal("Bob"))
	g.Expect(resultingVotes).To(Equal(map[string]bool{"Alice": true, "Bob": true, "Charlie": false, "Dan": false, "Edith": false}))
	g.Expect(newGame.teamVotes).To(BeNil())
	g.Expect(newGame.currentTeam).To(BeNil())
}

func Test_ApproveRejectTeam_ShouldMoveToGameOverIfVoteFailed5TimesInARow(t *testing.T) {
	g := NewWithT(t)
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.RejectTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, _, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.Leader()).To(Equal("Bob"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.RejectTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, _, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.Leader()).To(Equal("Charlie"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.RejectTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, _, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.Leader()).To(Equal("Dan"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.RejectTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, _, _ = newGame.RejectTeamBy("Edith")

	g.Expect(newGame.Leader()).To(Equal("Edith"))
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.RejectTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, resultingVotes, err := newGame.RejectTeamBy("Edith")

	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(GameOver))
	g.Expect(newGame.VoteFailures()).To(Equal(5))
	g.Expect(newGame.Leader()).To(Equal("Edith"))
	g.Expect(newGame.Winner()).To(Equal(Spy))
	g.Expect(resultingVotes).To(Equal(map[string]bool{"Alice": true, "Bob": true, "Charlie": false, "Dan": false, "Edith": false}))
	g.Expect(newGame.teamVotes).To(BeNil())
	g.Expect(newGame.currentTeam).To(BeNil())
}

func Test_ApproveRejectTeam_VoteFailureShouldResetAfterASuccessfulVote(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.RejectTeamBy("Charlie")
	newGame, _, _ = newGame.RejectTeamBy("Dan")
	newGame, _, _ = newGame.RejectTeamBy("Edith")

	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, resultingVotes, err := newGame.ApproveTeamBy("Edith")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(ConductingMission))
	g.Expect(newGame.VoteFailures()).To(Equal(0))
	g.Expect(resultingVotes).To(Equal(map[string]bool{"Alice": true, "Bob": true, "Charlie": true, "Dan": true, "Edith": true}))
	g.Expect(newGame.teamVotes).To(BeNil())
}

func Test_SucceedMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, outcomes, err := newGame.SucceedMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.missionOutcomes).To(Equal(votes(map[string]bool{"Alice": true})))
	g.Expect(outcomes).To(Equal(map[string]bool{"Alice": true}))
}

func Test_SucceedMission_ShouldFailIfPersonIsntOnTeam(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	_, outcomes, err := newGame.SucceedMissionBy("Charlie")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
	g.Expect(outcomes).To(BeNil())
}

func Test_SucceedMission_ShouldFailIfPersonAlreadyWorkedOnMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	_, outcomes, err := newGame.SucceedMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerHasAlreadyVoted))
	g.Expect(outcomes).To(BeNil())
}

func Test_SucceedMission_ShouldFailIfStateIsntConductingMission(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	_, outcomes, err := newGame.SucceedMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
	g.Expect(outcomes).To(BeNil())
}

func Test_FailMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, outcomes, err := newGame.FailMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.missionOutcomes).To(Equal(votes(map[string]bool{"Alice": false})))
	g.Expect(outcomes).To(Equal(map[string]bool{"Alice": false}))
}

func Test_FailMission_ShouldFailIfPersonIsntOnTeam(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	_, outcomes, err := newGame.FailMissionBy("Charlie")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
	g.Expect(outcomes).To(BeNil())
}

func Test_FailMission_ShouldFailIfPersonAlreadyWorkedOnMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _, _ = newGame.FailMissionBy("Alice")
	_, outcomes, err := newGame.FailMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerHasAlreadyVoted))
	g.Expect(outcomes).To(BeNil())
}

func Test_FailMission_ShouldFailIfStateIsntConductingMission(t *testing.T) {
	newGame := createNewlyVotingOnTeamGame()
	_, outcomes, err := newGame.FailMissionBy("Alice")

	g := NewWithT(t)
	g.Expect(err).To(MatchError(errInvalidStateForAction))
	g.Expect(outcomes).To(BeNil())
}

func Test_SucceedFailMission_ShouldMoveToSelectingTeamWhenEveryoneWorkedOnTheMission(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _, _ = newGame.FailMissionBy("Alice")
	newGame, outcomes, err := newGame.FailMissionBy("Bob")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(SelectingTeam))
	g.Expect(newGame.CurrentMission()).To(Equal(Second))
	g.Expect(newGame.missionOutcomes).To(BeNil())
	g.Expect(outcomes).To(Equal(map[string]bool{"Alice": false, "Bob": false}))
	g.Expect(newGame.GetMissionResults()).To(Equal(map[Mission]bool{First: false}))
	g.Expect(newGame.Leader()).To(Equal("Bob"))
	g.Expect(newGame.currentTeam).To(BeNil())
}

func Test_SucceedFailMission_ShouldMoveToGameOverIfThirdSuccess(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")

	// Second turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")
	newGame, _, _ = newGame.SucceedMissionBy("Charlie")

	// Third turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, outcomes, err := newGame.SucceedMissionBy("Bob")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(GameOver))
	g.Expect(newGame.CurrentMission()).To(Equal(Third))
	g.Expect(newGame.missionOutcomes).To(BeNil())
	g.Expect(newGame.GetMissionResults()).To(Equal(map[Mission]bool{First: true, Second: true, Third: true}))
	g.Expect(newGame.Leader()).To(Equal("Charlie"))
	g.Expect(newGame.Winner()).To(Equal(Resistance))
	g.Expect(outcomes).To(Equal(map[string]bool{"Alice": true, "Bob": true}))
}

func Test_SucceedFailMission_ShouldMoveToGameOverIfThirdFailure(t *testing.T) {
	newGame := createNewlyConductingMissionGame()
	newGame, _, _ = newGame.FailMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")

	// Second turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.FailMissionBy("Bob")
	newGame, _, _ = newGame.SucceedMissionBy("Charlie")

	// Third turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, outcomes, err := newGame.FailMissionBy("Bob")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.State()).To(Equal(GameOver))
	g.Expect(newGame.CurrentMission()).To(Equal(Third))
	g.Expect(newGame.missionOutcomes).To(BeNil())
	g.Expect(newGame.GetMissionResults()).To(Equal(map[Mission]bool{First: false, Second: false, Third: false}))
	g.Expect(newGame.Leader()).To(Equal("Charlie"))
	g.Expect(newGame.Winner()).To(Equal(Spy))
	g.Expect(outcomes).To(Equal(map[string]bool{"Alice": true, "Bob": false}))
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
	newGame, _, _, _ = newGame.Start(spiesFirstGenerator{})

	// First turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.ApproveTeamBy("Fred")
	newGame, _, _ = newGame.ApproveTeamBy("Gordon")
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")

	// Second turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.ApproveTeamBy("Fred")
	newGame, _, _ = newGame.ApproveTeamBy("Gordon")
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")
	newGame, _, _ = newGame.FailMissionBy("Charlie")

	// Third turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.ApproveTeamBy("Fred")
	newGame, _, _ = newGame.ApproveTeamBy("Gordon")
	newGame, _, _ = newGame.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")
	newGame, _, _ = newGame.FailMissionBy("Charlie")

	// Fourth turn
	newGame, _ = newGame.LeaderSelectsMember("Alice")
	newGame, _ = newGame.LeaderSelectsMember("Bob")
	newGame, _ = newGame.LeaderSelectsMember("Charlie")
	newGame, _ = newGame.LeaderSelectsMember("Dan")
	newGame, _ = newGame.LeaderConfirmsTeamSelection()
	newGame, _, _ = newGame.ApproveTeamBy("Alice")
	newGame, _, _ = newGame.ApproveTeamBy("Bob")
	newGame, _, _ = newGame.ApproveTeamBy("Charlie")
	newGame, _, _ = newGame.ApproveTeamBy("Dan")
	newGame, _, _ = newGame.ApproveTeamBy("Edith")
	newGame, _, _ = newGame.ApproveTeamBy("Fred")
	conductingFourthMission, _, _ := newGame.ApproveTeamBy("Gordon")
	newGame, _, _ = conductingFourthMission.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")
	newGame, _, _ = newGame.SucceedMissionBy("Charlie")
	newGame, outcomes, _ := newGame.FailMissionBy("Dan")

	g := NewWithT(t)
	g.Expect(newGame.GetMissionResults()).To(Equal(
		map[Mission]bool{
			First:  true,
			Second: false,
			Third:  false,
			Fourth: true,
		}))
	g.Expect(outcomes).To(Equal(map[string]bool{"Alice": true, "Bob": true, "Charlie": true, "Dan": false}))

	newGame, _, _ = conductingFourthMission.SucceedMissionBy("Alice")
	newGame, _, _ = newGame.SucceedMissionBy("Bob")
	newGame, _, _ = newGame.FailMissionBy("Charlie")
	newGame, outcomes, _ = newGame.FailMissionBy("Dan")
	g.Expect(newGame.GetMissionResults()).To(Equal(
		map[Mission]bool{
			First:  true,
			Second: false,
			Third:  false,
			Fourth: false,
		}))
	g.Expect(outcomes).To(Equal(map[string]bool{"Alice": true, "Bob": true, "Charlie": false, "Dan": false}))
}

func Test_Winner_ReturnsEmptyStringIfNotGameOver(t *testing.T) {
	newGame := createNewlyStartedGame()
	g := NewWithT(t)
	g.Expect(newGame.Winner()).To(Equal(Allegiance("")))
}
