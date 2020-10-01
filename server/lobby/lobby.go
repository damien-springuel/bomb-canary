package lobby

import (
	"fmt"

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

type lobbyServer struct {
	partyService partyService
}

func Register(engine *gin.Engine, partyService partyService) {

	lobbyServer := lobbyServer{
		partyService: partyService,
	}

	lobbyGroup := engine.Group("/lobby")
	lobbyGroup.GET("/create-party", lobbyServer.createParty)
}

func (l lobbyServer) createParty(c *gin.Context) {
	var req createPartyRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, fmt.Sprintf("%v", err))
		return
	}
	newCode := l.partyService.CreateParty()
	l.partyService.JoinParty(newCode, req.Name)
	c.JSON(200, createPartyResponse{Code: newCode})
}
