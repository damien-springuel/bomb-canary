package gamehub

import (
	"sync"

	"github.com/damien-springuel/bomb-canary/server/gamerules"
)

type games struct {
	mut         *sync.RWMutex
	gamesByCode map[string]gamerules.Game
}

func newGames() games {
	return games{
		mut:         &sync.RWMutex{},
		gamesByCode: make(map[string]gamerules.Game),
	}
}

func (g games) create(code string) {
	g.mut.Lock()
	g.gamesByCode[code] = gamerules.NewGame()
	g.mut.Unlock()
}

func (g games) get(code string) (gamerules.Game, bool) {
	g.mut.RLock()
	game, exists := g.gamesByCode[code]
	g.mut.RUnlock()
	return game, exists
}

func (g games) set(code string, game gamerules.Game) {
	g.mut.Lock()
	g.gamesByCode[code] = game
	g.mut.Unlock()
}
