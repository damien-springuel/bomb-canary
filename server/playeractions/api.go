package playeractions

import (
	"github.com/gin-gonic/gin"
)

const (
	partyCodeKey  = "partyCode"
	playerNameKey = "playerName"
)

type sessionGetter interface {
	Get(session string) (code string, name string, err error)
}

type actionBroker interface {
	StartGame(code string)
}

type playerActionServer struct {
	sessionGetter sessionGetter
	actionBroker  actionBroker
}

func Register(engine *gin.Engine, sessionGetter sessionGetter, actionBroker actionBroker) {

	playerActionServer := playerActionServer{
		sessionGetter: sessionGetter,
		actionBroker:  actionBroker,
	}

	actions := engine.Group("/actions")
	actions.Use(playerActionServer.checkSession)
	actions.GET("/start-game", playerActionServer.startGame)
}

func (p playerActionServer) checkSession(c *gin.Context) {
	session, err := c.Cookie("session")
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	partyCode, playerName, err := p.sessionGetter.Get(session)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

	setCodeAndNameToContext(c, partyCode, playerName)

	c.Next()
}

func setCodeAndNameToContext(c *gin.Context, code, name string) {
	c.Set(partyCodeKey, code)
	c.Set(playerNameKey, name)
}

func getCodeAndNameFromContext(c *gin.Context) (code, name string) {
	code = c.GetString(partyCodeKey)
	name = c.GetString(playerNameKey)
	return
}

func (p playerActionServer) startGame(c *gin.Context) {
	code, _ := getCodeAndNameFromContext(c)
	p.actionBroker.StartGame(code)

	c.JSON(200, gin.H{})
}
