package playeractions

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	partyCodeKey  = "partyCode"
	playerNameKey = "playerName"
)

type leaderSelectionRequest struct {
	Member string `json:"member"`
}

type sessionGetter interface {
	Get(session string) (code string, name string, err error)
}

type actionBroker interface {
	StartGame(code string)
	LeaderSelectsMember(code string, leader string, member string)
	LeaderDeselectsMember(code string, leader string, member string)
	LeaderConfirmsTeam(code string, leader string)
	ApproveTeam(code string, player string)
	RejectTeam(code string, player string)
	SucceedMission(code string, player string)
	FailMission(code string, player string)
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
	actions.POST("/start-game", playerActionServer.startGame)
	actions.POST("/leader-selects-member", playerActionServer.leaderSelectsMember)
	actions.POST("/leader-deselects-member", playerActionServer.leaderDeselectsMember)
	actions.POST("/leader-confirms-team", playerActionServer.leaderConfirmsTeam)
	actions.POST("/approve-team", playerActionServer.approveTeam)
	actions.POST("/reject-team", playerActionServer.rejectTeam)
	actions.POST("/succeed-mission", playerActionServer.succeedMission)
	actions.POST("/fail-mission", playerActionServer.failMission)
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

func (p playerActionServer) leaderSelectsMember(c *gin.Context) {
	var req leaderSelectionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("can't bind json: %v", err)})
		return
	}

	if req.Member == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "member is required"})
		return
	}

	code, name := getCodeAndNameFromContext(c)
	p.actionBroker.LeaderSelectsMember(code, name, req.Member)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) leaderDeselectsMember(c *gin.Context) {
	var req leaderSelectionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("can't bind json: %v", err)})
		return
	}

	if req.Member == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "member is required"})
		return
	}

	code, name := getCodeAndNameFromContext(c)
	p.actionBroker.LeaderDeselectsMember(code, name, req.Member)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) leaderConfirmsTeam(c *gin.Context) {
	code, name := getCodeAndNameFromContext(c)
	p.actionBroker.LeaderConfirmsTeam(code, name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) approveTeam(c *gin.Context) {
	code, name := getCodeAndNameFromContext(c)
	p.actionBroker.ApproveTeam(code, name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) rejectTeam(c *gin.Context) {
	code, name := getCodeAndNameFromContext(c)
	p.actionBroker.RejectTeam(code, name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) succeedMission(c *gin.Context) {
	code, name := getCodeAndNameFromContext(c)
	p.actionBroker.SucceedMission(code, name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) failMission(c *gin.Context) {
	code, name := getCodeAndNameFromContext(c)
	p.actionBroker.FailMission(code, name)

	c.JSON(200, gin.H{})
}
