package playeractions

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	playerNameKey = "playerName"
)

type leaderSelectionRequest struct {
	Member string `json:"member"`
}

type sessionGetter interface {
	Get(session string) (name string, err error)
}

type actionBroker interface {
	StartGame()
	LeaderSelectsMember(leader string, member string)
	LeaderDeselectsMember(leader string, member string)
	LeaderConfirmsTeam(leader string)
	ApproveTeam(player string)
	RejectTeam(player string)
	SucceedMission(player string)
	FailMission(player string)
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

	playerName, err := p.sessionGetter.Get(session)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

	setCodeAndNameToContext(c, playerName)

	c.Next()
}

func setCodeAndNameToContext(c *gin.Context, name string) {
	c.Set(playerNameKey, name)
}

func getNameFromContext(c *gin.Context) (name string) {
	name = c.GetString(playerNameKey)
	return
}

func (p playerActionServer) startGame(c *gin.Context) {
	p.actionBroker.StartGame()
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

	name := getNameFromContext(c)
	p.actionBroker.LeaderSelectsMember(name, req.Member)

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

	name := getNameFromContext(c)
	p.actionBroker.LeaderDeselectsMember(name, req.Member)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) leaderConfirmsTeam(c *gin.Context) {
	name := getNameFromContext(c)
	p.actionBroker.LeaderConfirmsTeam(name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) approveTeam(c *gin.Context) {
	name := getNameFromContext(c)
	p.actionBroker.ApproveTeam(name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) rejectTeam(c *gin.Context) {
	name := getNameFromContext(c)
	p.actionBroker.RejectTeam(name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) succeedMission(c *gin.Context) {
	name := getNameFromContext(c)
	p.actionBroker.SucceedMission(name)

	c.JSON(200, gin.H{})
}

func (p playerActionServer) failMission(c *gin.Context) {
	name := getNameFromContext(c)
	p.actionBroker.FailMission(name)

	c.JSON(200, gin.H{})
}
