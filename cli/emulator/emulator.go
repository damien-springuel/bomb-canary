package emulator

import (
	"context"
	"fmt"
	"image"
	"strconv"

	"github.com/damien-springuel/bomb-canary/cli/bcclient"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type pageType string
type emKey string

const (
	presets pageType = "presets"
	actions pageType = "actions"
	people  pageType = "people"
)

type Emulator struct {
	ctx   context.Context
	list  *widgets.List
	pages map[pageType]page
}

type choice struct {
	description string
	action      func(ctx context.Context) context.Context
}
type page struct {
	title string
	rows  []choice
}

func getValueFromContext(ctx context.Context, key string) any {
	return ctx.Value(emKey(key))
}

func setValueToContext(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, emKey(key), value)
}

func getCurrentPageFromContext(ctx context.Context) pageType {
	return getPageFromContext(ctx, "currentPage")
}

func setCurrentPageToContext(ctx context.Context, currentPage pageType) context.Context {
	return setValueToContext(ctx, "currentPage", currentPage)
}

func getNextPageFromContext(ctx context.Context) pageType {
	return getPageFromContext(ctx, "nextPage")
}

func setNextPageToContext(ctx context.Context, nextPage pageType) context.Context {
	return setValueToContext(ctx, "nextPage", nextPage)
}
func setNilNextPageToContext(ctx context.Context) context.Context {
	return setNextPageToContext(ctx, "")
}

func getPageFromContext(ctx context.Context, page string) pageType {
	p, ok := getValueFromContext(ctx, page).(pageType)
	if !ok {
		return pageType("")
	}
	return p
}

func getActionFromContext(ctx context.Context) func(context.Context) context.Context {
	return getValueFromContext(ctx, "action").(func(context.Context) context.Context)
}

func setActionToContext(ctx context.Context, action func(context.Context) context.Context) context.Context {
	return setValueToContext(ctx, "action", action)
}

func getActionDescFromContext(ctx context.Context) string {
	desc, ok := getValueFromContext(ctx, "actionDesc").(string)
	if !ok {
		return ""
	}
	return desc
}
func setActionDescToContext(ctx context.Context, actionDesc string) context.Context {
	return setValueToContext(ctx, "actionDesc", actionDesc)
}

func setEmptyActionDescToContext(ctx context.Context) context.Context {
	return setValueToContext(ctx, "actionDesc", "")
}

func getSessionFromContext(ctx context.Context, name string) string {
	return getValueFromContext(ctx, "session"+name).(string)
}

func setSessionToContext(ctx context.Context, name string, session string) context.Context {
	return setValueToContext(ctx, "session"+name, session)
}

func getNameFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, "name").(string)
}

func setNameToContext(ctx context.Context, name string) context.Context {
	return setValueToContext(ctx, "name", name)
}

func getLeaderFromContext(ctx context.Context) string {
	return getValueFromContext(ctx, "leader").(string)
}

func setLeaderToContext(ctx context.Context, leader string) context.Context {
	return setValueToContext(ctx, "leader", leader)
}

func setCodeToContext(ctx context.Context, code string) context.Context {
	return setValueToContext(ctx, "code", code)
}

func createSetNameAction(name string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = setNameToContext(ctx, name)
		action := getActionFromContext(ctx)
		ctx = action(ctx)
		nextPage := getNextPageFromContext(ctx)
		if nextPage == "" {
			nextPage = actions
			ctx = setEmptyActionDescToContext(ctx)
		}
		return setCurrentPageToContext(ctx, nextPage)
	}
}

func createActionWithName(a func(ctx context.Context) context.Context, actionDesc string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = setActionToContext(ctx, a)
		ctx = setActionDescToContext(ctx, actionDesc)
		return setCurrentPageToContext(ctx, people)
	}
}

func blankPreset(ctx context.Context) context.Context {
	return setCurrentPageToContext(ctx, actions)
}

func createdAndJoined5Players(ctx context.Context) context.Context {
	ctx = setCurrentPageToContext(ctx, actions)
	code, session := bcclient.CreateGame("Alice")
	ctx = setCodeToContext(ctx, code)
	ctx = setSessionToContext(ctx, "Alice", session)
	session = bcclient.JoinGame(code, "Bob")
	ctx = setSessionToContext(ctx, "Bob", session)
	session = bcclient.JoinGame(code, "Charlie")
	ctx = setSessionToContext(ctx, "Charlie", session)
	session = bcclient.JoinGame(code, "Dan")
	ctx = setSessionToContext(ctx, "Dan", session)
	session = bcclient.JoinGame(code, "Edith")
	ctx = setSessionToContext(ctx, "Edith", session)
	return ctx
}

func createdAndJoined10Players(ctx context.Context) context.Context {
	ctx = setCurrentPageToContext(ctx, actions)
	code, session := bcclient.CreateGame("Alice")
	ctx = setCodeToContext(ctx, code)
	ctx = setSessionToContext(ctx, "Alice", session)
	session = bcclient.JoinGame(code, "Bob")
	ctx = setSessionToContext(ctx, "Bob", session)
	session = bcclient.JoinGame(code, "Charlie")
	ctx = setSessionToContext(ctx, "Charlie", session)
	session = bcclient.JoinGame(code, "Dan")
	ctx = setSessionToContext(ctx, "Dan", session)
	session = bcclient.JoinGame(code, "Edith")
	ctx = setSessionToContext(ctx, "Edith", session)
	session = bcclient.JoinGame(code, "Frank")
	ctx = setSessionToContext(ctx, "Frank", session)
	session = bcclient.JoinGame(code, "Gus")
	ctx = setSessionToContext(ctx, "Gus", session)
	session = bcclient.JoinGame(code, "Henry")
	ctx = setSessionToContext(ctx, "Henry", session)
	session = bcclient.JoinGame(code, "Ian")
	ctx = setSessionToContext(ctx, "Ian", session)
	session = bcclient.JoinGame(code, "Jay")
	ctx = setSessionToContext(ctx, "Jay", session)
	return ctx
}

func started5PlayerGame(ctx context.Context) context.Context {
	ctx = createdAndJoined5Players(ctx)
	session := getSessionFromContext(ctx, "Alice")
	bcclient.StartGame(session)
	return ctx
}

func started10PlayerGame(ctx context.Context) context.Context {
	ctx = createdAndJoined10Players(ctx)
	session := getSessionFromContext(ctx, "Alice")
	bcclient.StartGame(session)
	return ctx
}

func fivePlayerGameFirstVote(ctx context.Context) context.Context {
	ctx = started5PlayerGame(ctx)
	session := getSessionFromContext(ctx, "Alice")
	bcclient.LeaderSelectsMember(session, "Alice")
	bcclient.LeaderSelectsMember(session, "Bob")
	bcclient.LeaderConfirmsTeam(session)
	return ctx
}

func fivePlayerGameFirstMission(ctx context.Context) context.Context {
	ctx = fivePlayerGameFirstVote(ctx)
	bcclient.ApproveTeam(getSessionFromContext(ctx, "Alice"))
	bcclient.ApproveTeam(getSessionFromContext(ctx, "Bob"))
	bcclient.ApproveTeam(getSessionFromContext(ctx, "Charlie"))
	bcclient.ApproveTeam(getSessionFromContext(ctx, "Dan"))
	bcclient.ApproveTeam(getSessionFromContext(ctx, "Edith"))
	return ctx
}

func fivePlayerGameOneSuccessOneFailure(ctx context.Context) context.Context {
	ctx = started5PlayerGame(ctx)
	alice := getSessionFromContext(ctx, "Alice")
	bob := getSessionFromContext(ctx, "Bob")
	charlie := getSessionFromContext(ctx, "Charlie")
	dan := getSessionFromContext(ctx, "Dan")
	edith := getSessionFromContext(ctx, "Edith")

	bcclient.LeaderSelectsMember(alice, "Alice")
	bcclient.LeaderSelectsMember(alice, "Bob")
	bcclient.LeaderConfirmsTeam(alice)
	bcclient.ApproveTeam(alice)
	bcclient.ApproveTeam(bob)
	bcclient.RejectTeam(charlie)
	bcclient.RejectTeam(dan)
	bcclient.RejectTeam(edith)

	bcclient.LeaderSelectsMember(bob, "Charlie")
	bcclient.LeaderSelectsMember(bob, "Bob")
	bcclient.LeaderConfirmsTeam(bob)
	bcclient.ApproveTeam(alice)
	bcclient.ApproveTeam(bob)
	bcclient.ApproveTeam(charlie)
	bcclient.RejectTeam(dan)
	bcclient.RejectTeam(edith)
	bcclient.SucceedMission(charlie)
	bcclient.SucceedMission(bob)

	bcclient.LeaderSelectsMember(charlie, "Alice")
	bcclient.LeaderSelectsMember(charlie, "Bob")
	bcclient.LeaderSelectsMember(charlie, "Charlie")
	bcclient.LeaderConfirmsTeam(charlie)
	bcclient.RejectTeam(alice)
	bcclient.ApproveTeam(bob)
	bcclient.RejectTeam(charlie)
	bcclient.ApproveTeam(dan)
	bcclient.RejectTeam(edith)

	bcclient.LeaderSelectsMember(dan, "Dan")
	bcclient.LeaderSelectsMember(dan, "Edith")
	bcclient.LeaderSelectsMember(dan, "Charlie")
	bcclient.LeaderConfirmsTeam(dan)
	bcclient.ApproveTeam(alice)
	bcclient.RejectTeam(bob)
	bcclient.RejectTeam(charlie)
	bcclient.ApproveTeam(dan)
	bcclient.ApproveTeam(edith)
	bcclient.SucceedMission(dan)
	bcclient.FailMission(edith)
	bcclient.SucceedMission(charlie)
	return ctx
}

func tenPlayerGameFourTeamVoteFailures(ctx context.Context) context.Context {
	ctx = started10PlayerGame(ctx)
	alice := getSessionFromContext(ctx, "Alice")
	bob := getSessionFromContext(ctx, "Bob")
	charlie := getSessionFromContext(ctx, "Charlie")
	dan := getSessionFromContext(ctx, "Dan")
	edith := getSessionFromContext(ctx, "Edith")
	frank := getSessionFromContext(ctx, "Frank")
	gus := getSessionFromContext(ctx, "Gus")
	henry := getSessionFromContext(ctx, "Henry")
	ian := getSessionFromContext(ctx, "Ian")
	jay := getSessionFromContext(ctx, "Jay")

	bcclient.LeaderSelectsMember(alice, "Alice")
	bcclient.LeaderSelectsMember(alice, "Bob")
	bcclient.LeaderSelectsMember(alice, "Charlie")
	bcclient.LeaderConfirmsTeam(alice)
	bcclient.RejectTeam(alice)
	bcclient.RejectTeam(bob)
	bcclient.RejectTeam(charlie)
	bcclient.RejectTeam(dan)
	bcclient.RejectTeam(edith)
	bcclient.RejectTeam(frank)
	bcclient.RejectTeam(gus)
	bcclient.RejectTeam(henry)
	bcclient.RejectTeam(ian)
	bcclient.RejectTeam(jay)

	bcclient.LeaderSelectsMember(bob, "Dan")
	bcclient.LeaderSelectsMember(bob, "Gus")
	bcclient.LeaderSelectsMember(bob, "Jay")
	bcclient.LeaderConfirmsTeam(bob)
	bcclient.RejectTeam(alice)
	bcclient.RejectTeam(bob)
	bcclient.RejectTeam(charlie)
	bcclient.ApproveTeam(dan)
	bcclient.RejectTeam(edith)
	bcclient.RejectTeam(frank)
	bcclient.RejectTeam(gus)
	bcclient.RejectTeam(henry)
	bcclient.ApproveTeam(ian)
	bcclient.ApproveTeam(jay)

	bcclient.LeaderSelectsMember(charlie, "Edith")
	bcclient.LeaderSelectsMember(charlie, "Alice")
	bcclient.LeaderSelectsMember(charlie, "Frank")
	bcclient.LeaderConfirmsTeam(charlie)
	bcclient.ApproveTeam(alice)
	bcclient.RejectTeam(bob)
	bcclient.ApproveTeam(charlie)
	bcclient.RejectTeam(dan)
	bcclient.ApproveTeam(edith)
	bcclient.ApproveTeam(frank)
	bcclient.RejectTeam(gus)
	bcclient.ApproveTeam(henry)
	bcclient.RejectTeam(ian)
	bcclient.RejectTeam(jay)

	bcclient.LeaderSelectsMember(dan, "Ian")
	bcclient.LeaderSelectsMember(dan, "Henry")
	bcclient.LeaderSelectsMember(dan, "Gus")
	bcclient.LeaderConfirmsTeam(dan)
	bcclient.ApproveTeam(alice)
	bcclient.RejectTeam(bob)
	bcclient.ApproveTeam(charlie)
	bcclient.RejectTeam(dan)
	bcclient.RejectTeam(edith)
	bcclient.RejectTeam(frank)
	bcclient.RejectTeam(gus)
	bcclient.RejectTeam(henry)
	bcclient.RejectTeam(ian)
	bcclient.ApproveTeam(jay)

	return ctx
}

func New() *Emulator {
	presetsPage := page{
		title: "Presets",
		rows: []choice{
			{
				description: "Blank",
				action:      blankPreset,
			},
			{
				description: "Created and joined 5 players",
				action:      createdAndJoined5Players,
			},
			{
				description: "Started 5-player game",
				action:      started5PlayerGame,
			},
			{
				description: "5-player game first vote",
				action:      fivePlayerGameFirstVote,
			},
			{
				description: "5-player game first mission",
				action:      fivePlayerGameFirstMission,
			},
			{
				description: "5-player 1 success, 1 failure",
				action:      fivePlayerGameOneSuccessOneFailure,
			},
			{
				description: "Started 10-player game",
				action:      started10PlayerGame,
			},
			{
				description: "10 player game with 4 team vote failures",
				action:      tenPlayerGameFourTeamVoteFailures,
			},
		},
	}
	peoplePage := page{
		title: "Who ",
		rows: []choice{
			{
				description: "Alice",
				action:      createSetNameAction("Alice"),
			},
			{
				description: "Bob",
				action:      createSetNameAction("Bob"),
			},
			{
				description: "Charlie",
				action:      createSetNameAction("Charlie"),
			},
			{
				description: "Dan",
				action:      createSetNameAction("Dan"),
			},
			{
				description: "Edith",
				action:      createSetNameAction("Edith"),
			},
		},
	}
	actionPage := page{
		title: "Action",
		rows: []choice{
			{
				description: "Start Game",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					session := getSessionFromContext(ctx, name)
					bcclient.StartGame(session)
					return ctx
				}, "is starting the game?"),
			},
			{
				description: "Leader Selects Member",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					ctx = setLeaderToContext(ctx, name)
					ctx = setNextPageToContext(ctx, people)
					ctx = setActionToContext(ctx, func(ctx context.Context) context.Context {
						leader := getLeaderFromContext(ctx)
						session := getSessionFromContext(ctx, leader)
						name := getNameFromContext(ctx)
						bcclient.LeaderSelectsMember(session, name)
						return setNilNextPageToContext(ctx)
					})
					ctx = setActionDescToContext(ctx, "is being selected?")
					return setCurrentPageToContext(ctx, people)
				}, "is the leader?"),
			},
			{
				description: "Leader Deselects Member",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					ctx = setLeaderToContext(ctx, name)
					ctx = setNextPageToContext(ctx, people)
					ctx = setActionToContext(ctx, func(ctx context.Context) context.Context {
						leader := getLeaderFromContext(ctx)
						session := getSessionFromContext(ctx, leader)
						name := getNameFromContext(ctx)
						bcclient.LeaderDeselectsMember(session, name)
						return setNilNextPageToContext(ctx)
					})
					ctx = setActionDescToContext(ctx, "is being deselected")
					return setCurrentPageToContext(ctx, people)
				}, "is the leader?"),
			},
			{
				description: "Leader Confirms Team",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					session := getSessionFromContext(ctx, name)
					bcclient.LeaderConfirmsTeam(session)
					return ctx
				}, "is confirming the team?"),
			},
			{
				description: "Approve Team",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					session := getSessionFromContext(ctx, name)
					bcclient.ApproveTeam(session)
					return ctx
				}, "is approving the team?"),
			},
			{
				description: "Reject Team",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					session := getSessionFromContext(ctx, name)
					bcclient.RejectTeam(session)
					return ctx
				}, "is rejecting the team?"),
			},
			{
				description: "Succeed Mission",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					session := getSessionFromContext(ctx, name)
					bcclient.SucceedMission(session)
					return ctx
				}, "is succeeding the mission?"),
			},
			{
				description: "Fail Mission",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := getNameFromContext(ctx)
					session := getSessionFromContext(ctx, name)
					bcclient.FailMission(session)
					return ctx
				}, "is failing the mission?"),
			},
		},
	}

	list := widgets.NewList()
	return &Emulator{
		list: list,
		pages: map[pageType]page{
			actions: actionPage,
			people:  peoplePage,
			presets: presetsPage,
		},
		ctx: setCurrentPageToContext(context.Background(), presets),
	}
}

func (e *Emulator) updateList() {
	p := e.getCurrentPage()
	titleDesc := getActionDescFromContext(e.ctx)
	e.list.Title = fmt.Sprintf("%s%s", p.title, titleDesc)
	rows := make([]string, len(p.rows))
	for i := range rows {
		rows[i] = fmt.Sprintf("%d. %s", i+1, p.rows[i].description)
	}
	e.list.Rows = rows
}

func (e *Emulator) getCurrentPage() page {
	return e.pages[getCurrentPageFromContext(e.ctx)]
}

func (e *Emulator) HandleUiEvent(event termui.Event) {
	index, err := strconv.ParseInt(event.ID, 10, 0)
	if err != nil {
		return
	}
	if int(index) > len(e.getCurrentPage().rows) {
		return
	}
	e.ctx = e.getCurrentPage().rows[index-1].action(e.ctx)
}

// Drawable functions
func (e *Emulator) GetRect() image.Rectangle {
	return e.list.GetRect()
}

func (e *Emulator) SetRect(x1, y1, x2, y2 int) {
	e.list.SetRect(x1, y1, x2, y2)
}

func (e *Emulator) Draw(buf *termui.Buffer) {
	e.updateList()
	e.list.Draw(buf)
}

func (e *Emulator) Lock() {
	e.list.Lock()
}

func (e *Emulator) Unlock() {
	e.list.Unlock()
}
