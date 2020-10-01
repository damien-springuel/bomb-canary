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

type joinPartyRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
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
	lobbyGroup.GET("/join-party", lobbyServer.joinParty)
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

	setSessionCookie(c, l.session.Create(newCode, req.Name))

	c.JSON(200, createPartyResponse{Code: newCode})
}

func (l lobbyServer) joinParty(c *gin.Context) {
	var req joinPartyRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("can't bind json: %v", err)})
		return
	}

	if req.Code == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "code is required"})
		return
	}

	if req.Name == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "name is required"})
		return
	}

	l.partyService.JoinParty(req.Code, req.Name)

	setSessionCookie(c, l.session.Create(req.Code, req.Name))

	c.JSON(200, gin.H{})
}

func setSessionCookie(c *gin.Context, session string) {
	c.SetCookie("session", session, int((time.Hour * 3).Seconds()), "/", "", false, true)
}
