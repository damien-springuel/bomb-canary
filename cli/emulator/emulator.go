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

func createSetNameAction(name string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, "name", name)
		action := ctx.Value("action").(func(context.Context) context.Context)
		ctx = action(ctx)
		nextPage := ctx.Value("nextPage")
		if nextPage == nil {
			nextPage = actions
			ctx = context.WithValue(ctx, "actionDesc", "")
		}
		return context.WithValue(ctx, "currentPage", nextPage.(pageType))
	}
}

func createActionWithName(a func(ctx context.Context) context.Context, actionDesc string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, "action", a)
		ctx = context.WithValue(ctx, "actionDesc", actionDesc)
		return context.WithValue(ctx, "currentPage", people)
	}
}

func blankPreset(ctx context.Context) context.Context {
	return context.WithValue(ctx, "currentPage", actions)
}

func createdAndJoined5Players(ctx context.Context) context.Context {
	ctx = blankPreset(ctx)
	ctx = context.WithValue(ctx, "currentPage", actions)
	code, session := bcclient.CreateGame("Alice")
	ctx = context.WithValue(ctx, "code", code)
	ctx = context.WithValue(ctx, "sessionAlice", session)
	session = bcclient.JoinGame(code, "Bob")
	ctx = context.WithValue(ctx, "sessionBob", session)
	session = bcclient.JoinGame(code, "Charlie")
	ctx = context.WithValue(ctx, "sessionCharlie", session)
	session = bcclient.JoinGame(code, "Dan")
	ctx = context.WithValue(ctx, "sessionDan", session)
	session = bcclient.JoinGame(code, "Edith")
	ctx = context.WithValue(ctx, "sessionEdith", session)
	return ctx
}

func started5PlayerGame(ctx context.Context) context.Context {
	ctx = createdAndJoined5Players(ctx)
	session := ctx.Value("sessionAlice").(string)
	bcclient.StartGame(session)
	return ctx
}

func fivePlayerGameFirstVote(ctx context.Context) context.Context {
	ctx = started5PlayerGame(ctx)
	session := ctx.Value("sessionAlice").(string)
	bcclient.LeaderSelectsMember(session, "Alice")
	bcclient.LeaderSelectsMember(session, "Bob")
	bcclient.LeaderConfirmsTeam(session)
	return ctx
}

func fivePlayerGameFirstMission(ctx context.Context) context.Context {
	ctx = fivePlayerGameFirstVote(ctx)
	bcclient.ApproveTeam(ctx.Value("sessionAlice").(string))
	bcclient.ApproveTeam(ctx.Value("sessionBob").(string))
	bcclient.ApproveTeam(ctx.Value("sessionCharlie").(string))
	bcclient.ApproveTeam(ctx.Value("sessionDan").(string))
	bcclient.ApproveTeam(ctx.Value("sessionEdith").(string))
	return ctx
}

func fivePlayerGameOneSuccessOneFailure(ctx context.Context) context.Context {
	ctx = started5PlayerGame(ctx)
	alice := ctx.Value("sessionAlice").(string)
	bob := ctx.Value("sessionBob").(string)
	charlie := ctx.Value("sessionCharlie").(string)
	dan := ctx.Value("sessionDan").(string)
	edith := ctx.Value("sessionEdith").(string)

	bcclient.LeaderSelectsMember(alice, "Alice")
	bcclient.LeaderSelectsMember(alice, "Bob")
	bcclient.LeaderConfirmsTeam(alice)
	bcclient.ApproveTeam(alice)
	bcclient.ApproveTeam(bob)
	bcclient.ApproveTeam(charlie)
	bcclient.ApproveTeam(dan)
	bcclient.ApproveTeam(edith)
	bcclient.SucceedMission(alice)
	bcclient.SucceedMission(bob)

	bcclient.LeaderSelectsMember(bob, "Alice")
	bcclient.LeaderSelectsMember(bob, "Bob")
	bcclient.LeaderSelectsMember(bob, "Charlie")
	bcclient.LeaderConfirmsTeam(bob)
	bcclient.ApproveTeam(alice)
	bcclient.ApproveTeam(bob)
	bcclient.ApproveTeam(charlie)
	bcclient.ApproveTeam(dan)
	bcclient.ApproveTeam(edith)
	bcclient.FailMission(alice)
	bcclient.FailMission(bob)
	bcclient.FailMission(charlie)
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
					name := ctx.Value("name").(string)
					session := ctx.Value("session" + name).(string)
					bcclient.StartGame(session)
					return ctx
				}, "is starting the game?"),
			},
			{
				description: "Leader Selects Member",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					ctx = context.WithValue(ctx, "leader", name)
					ctx = context.WithValue(ctx, "nextPage", people)
					ctx = context.WithValue(ctx, "action", func(ctx context.Context) context.Context {
						leader := ctx.Value("leader").(string)
						session := ctx.Value("session" + leader).(string)
						name := ctx.Value("name").(string)
						bcclient.LeaderSelectsMember(session, name)
						return context.WithValue(ctx, "nextPage", nil)
					})
					ctx = context.WithValue(ctx, "actionDesc", "is being selected?")
					return context.WithValue(ctx, "currentPage", people)
				}, "is the leader?"),
			},
			{
				description: "Leader Deselects Member",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					ctx = context.WithValue(ctx, "leader", name)
					ctx = context.WithValue(ctx, "nextPage", people)
					ctx = context.WithValue(ctx, "action", func(ctx context.Context) context.Context {
						leader := ctx.Value("leader").(string)
						session := ctx.Value("session" + leader).(string)
						name := ctx.Value("name").(string)
						bcclient.LeaderDeselectsMember(session, name)
						return context.WithValue(ctx, "nextPage", nil)
					})
					return context.WithValue(ctx, "actionDesc", "is being deselected?")
				}, "is the leader?"),
			},
			{
				description: "Leader Confirms Team",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					session := ctx.Value("session" + name).(string)
					bcclient.LeaderConfirmsTeam(session)
					return ctx
				}, "is confirming the team?"),
			},
			{
				description: "Approve Team",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					session := ctx.Value("session" + name).(string)
					bcclient.ApproveTeam(session)
					return ctx
				}, "is approving the team?"),
			},
			{
				description: "Reject Team",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					session := ctx.Value("session" + name).(string)
					bcclient.RejectTeam(session)
					return ctx
				}, "is rejecting the team?"),
			},
			{
				description: "Succeed Mission",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					session := ctx.Value("session" + name).(string)
					bcclient.SucceedMission(session)
					return ctx
				}, "is succeeding the mission?"),
			},
			{
				description: "Fail Mission",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					session := ctx.Value("session" + name).(string)
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
		ctx: context.WithValue(context.Background(), "currentPage", presets),
	}
}

func (e *Emulator) updateList() {
	p := e.getCurrentPage()
	titleDesc, ok := e.ctx.Value("actionDesc").(string)
	if !ok {
		titleDesc = ""
	}
	e.list.Title = fmt.Sprintf("%s%s", p.title, titleDesc)
	rows := make([]string, len(p.rows))
	for i := range rows {
		rows[i] = fmt.Sprintf("%d. %s", i+1, p.rows[i].description)
	}
	e.list.Rows = rows
}

func getPageFromContext(ctx context.Context) pageType {
	return ctx.Value("currentPage").(pageType)
}

func (e *Emulator) getCurrentPage() page {
	return e.pages[getPageFromContext(e.ctx)]
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
