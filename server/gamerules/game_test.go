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
