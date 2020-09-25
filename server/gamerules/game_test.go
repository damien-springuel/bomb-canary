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
	newGame = newGame.addPlayer("Alice")

	g := NewWithT(t)
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}

func Test_RemovePlayer(t *testing.T) {
	newGame := newGame("123")
	newGame = newGame.addPlayer("Alice")
	newGame, err := newGame.removePlayer("Alice")

	g := NewWithT(t)
	g.Expect(err).To(BeNil())
	g.Expect(newGame.players).To(Equal(players([]string{})))
}

func Test_RemovePlayer_PlayerNotFound(t *testing.T) {
	newGame := newGame("123")
	newGame = newGame.addPlayer("Alice")
	newGame, err := newGame.removePlayer("Bob")

	g := NewWithT(t)
	g.Expect(err).To(Equal(errPlayerNotFound))
	g.Expect(newGame.players).To(Equal(players([]string{"Alice"})))
}
