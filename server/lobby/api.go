package lobby

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	lobbyServer := lobbyServer{}

	lobbyGroup := engine.Group("/lobby")
	lobbyGroup.GET("/create-party", func(c *gin.Context) { handleCreateParty(c, lobbyServer) })
}

func handleCreateParty(c *gin.Context, lobbyServer lobbyServer) {
	var req createPartyRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, fmt.Sprintf("%v", err))
		return
	}
	resp, err := lobbyServer.createParty(req)
	if err != nil {
		c.JSON(500, fmt.Sprintf("%v", err))
		return
	}
	c.JSON(200, resp)
}
