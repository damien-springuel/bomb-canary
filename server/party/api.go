package party

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var errPartyDoesntExists = errors.New("party doesn't exist")

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

type partyBroker interface {
	CreateParty() string
	JoinParty(code string, name string) error
}

type sessionCreator interface {
	Create(code string, name string) string
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
	lobbyGroup.GET("/create", lobbyServer.createParty)
	lobbyGroup.GET("/join", lobbyServer.joinParty)
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

	newCode := l.partyBroker.CreateParty()
	_ = l.partyBroker.JoinParty(newCode, req.Name)

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

	err = l.partyBroker.JoinParty(req.Code, req.Name)
	if err != nil {
		if errors.Is(err, errPartyDoesntExists) {
			c.AbortWithStatusJSON(404, gin.H{"error": fmt.Sprintf("%v", err)})
			return
		}
		_ = c.AbortWithError(500, err)
		return
	}

	setSessionCookie(c, l.session.Create(req.Code, req.Name))

	c.JSON(200, gin.H{})
}

func setSessionCookie(c *gin.Context, session string) {
	c.SetCookie("session", session, int((time.Hour * 3).Seconds()), "/", "", false, true)
}
