package lobby

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type createPartyRequest struct {
	Name string `json:"name"`
}
type createPartyResponse struct {
	Code string `json:"code"`
}

type partyService interface {
	CreateParty() string
	JoinParty(code string, name string)
}

type sessionCreator interface {
	Create(code string, name string) string
}

type lobbyServer struct {
	partyService partyService
	session      sessionCreator
}

func Register(engine *gin.Engine, partyService partyService, session sessionCreator) {
	lobbyServer := lobbyServer{
		partyService: partyService,
		session:      session,
	}

	lobbyGroup := engine.Group("/lobby")
	lobbyGroup.GET("/create-party", lobbyServer.createParty)
}

func (l lobbyServer) createParty(c *gin.Context) {
	var req createPartyRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("can't bind json: %v", err)})
		return
	}

	if req.Name == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "name is required"})
		return
	}

	newCode := l.partyService.CreateParty()
	l.partyService.JoinParty(newCode, req.Name)

	session := l.session.Create(newCode, req.Name)
	c.SetCookie("session", session, int((time.Minute * 60).Seconds()), "/", "", false, true)

	c.JSON(200, createPartyResponse{Code: newCode})
}
