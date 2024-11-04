package party

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type joinPartyRequest struct {
	Name string `json:"name"`
}

type partyBroker interface {
	JoinParty(name string)
}

type sessionCreator interface {
	Create(name string) string
}

type lobbyServer struct {
	partyBroker partyBroker
	session     sessionCreator
}

func Register(engine *gin.Engine, partyBroker partyBroker, session sessionCreator) {
	lobbyServer := lobbyServer{
		partyBroker: partyBroker,
		session:     session,
	}

	lobbyGroup := engine.Group("/party")
	lobbyGroup.POST("/join", lobbyServer.joinParty)
}

func (l lobbyServer) joinParty(c *gin.Context) {
	var req joinPartyRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("can't bind json: %v", err)})
		return
	}

	if req.Name == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "name is required"})
		return
	}

	l.partyBroker.JoinParty(req.Name)
	setSessionCookie(c, l.session.Create(req.Name))

	c.JSON(200, gin.H{})
}

func setSessionCookie(c *gin.Context, session string) {
	c.SetCookie("session", session, int((5 * time.Hour).Seconds()), "/", "", false, true)
}
