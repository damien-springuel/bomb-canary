package gamerules

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_CreateGame(t *testing.T) {
	newGame := newGame("123")
	g := NewWithT(t)
	g.Expect(newGame).To(Equal(game{id: "123", state: notStarted}))
}

func Test_AddPlayer(t *testing.T) {
	newGame := newGame("123")
	newGame, _ = newGame.addPlayer("Alice")

	g := NewWithT(t)
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}

func Test_AddPlayer_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := newGame("123")
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

func Test_RemovePlayer(t *testing.T) {
	newGame := newGame("123")
	newGame, _ = newGame.addPlayer("Alice")
	newGame, err := newGame.removePlayer("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.players).To(Equal(players([]string{})))
}

func Test_RemovePlayer_ShouldErrorIfPlayerNotFound(t *testing.T) {
	newGame := newGame("123")
	newGame, _ = newGame.addPlayer("Alice")
	newGame, err := newGame.removePlayer("Bob")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}

func Test_RemovePlayer_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := newGame("123")
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
	newGame := newGame("123")
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")

	newGame, err := newGame.start()

	g := NewWithT(t)
	g.Expect(err).To(Equal(errNotEnoughPlayers))
}

func Test_StartGame(t *testing.T) {
	newGame := newGame("123")
	newGame, _ = newGame.addPlayer("Alice")
	newGame, _ = newGame.addPlayer("Bob")
	newGame, _ = newGame.addPlayer("Charlie")
	newGame, _ = newGame.addPlayer("Dan")
	newGame, _ = newGame.addPlayer("Edith")

	newGame, err := newGame.start()

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.state).To(Equal(selectingTeam))
}

func Test_StartGame_ShouldErrorIfGameHasStarted(t *testing.T) {
	newGame := newGame("123")
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
