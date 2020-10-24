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
	defer g.mut.Unlock()
	g.gamesByCode[code] = gamerules.NewGame()
}

func (g games) get(code string) (gamerules.Game, bool) {
	g.mut.RLock()
	defer g.mut.RUnlock()
	game, exists := g.gamesByCode[code]
	return game, exists
}

func (g games) set(code string, game gamerules.Game) {
	g.mut.Lock()
	defer g.mut.Unlock()
	g.gamesByCode[code] = game
}
