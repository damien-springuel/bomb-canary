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
		nameContext := context.WithValue(ctx, "name", name)
		action := ctx.Value("action").(func(context.Context) context.Context)
		resultCtx := action(nameContext)
		nextPage := resultCtx.Value("nextPage")
		actionDescCtx := resultCtx
		if nextPage == nil {
			nextPage = actions
			actionDescCtx = context.WithValue(resultCtx, "actionDesc", "")
		}
		return context.WithValue(actionDescCtx, "currentPage", nextPage.(pageType))
	}
}

func createActionWithName(a func(ctx context.Context) context.Context, actionDesc string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		actionCtx := context.WithValue(ctx, "action", a)
		actionDescCtx := context.WithValue(actionCtx, "actionDesc", actionDesc)
		return context.WithValue(actionDescCtx, "currentPage", people)
	}
}

func New() *Emulator {

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
				description: "Create Game",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					code, session := bcclient.CreateGame(name)
					codeCtx := context.WithValue(ctx, "code", code)
					return context.WithValue(codeCtx, "session"+name, session)
				}, "is creating the game?"),
			},
			{
				description: "Join Game",
				action: createActionWithName(func(ctx context.Context) context.Context {
					code := ctx.Value("code").(string)
					name := ctx.Value("name").(string)
					session := bcclient.JoinGame(code, name)
					return context.WithValue(ctx, "session"+name, session)
				}, "is joining the game?"),
			},
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
					leaderCtx := context.WithValue(ctx, "leader", name)
					nextPageCtx := context.WithValue(leaderCtx, "nextPage", people)
					actionCtx := context.WithValue(nextPageCtx, "action", func(ctx context.Context) context.Context {
						leader := ctx.Value("leader").(string)
						session := ctx.Value("session" + leader).(string)
						name := ctx.Value("name").(string)
						bcclient.LeaderSelectsMember(session, name)
						return context.WithValue(ctx, "nextPage", nil)
					})
					actionDescCtx := context.WithValue(actionCtx, "actionDesc", "is being selected?")
					return context.WithValue(actionDescCtx, "currentPage", people)
				}, "is the leader?"),
			},
			{
				description: "Leader Deselects Member",
				action: createActionWithName(func(ctx context.Context) context.Context {
					name := ctx.Value("name").(string)
					leaderCtx := context.WithValue(ctx, "leader", name)
					nextPageCtx := context.WithValue(leaderCtx, "nextPage", people)
					actionCtx := context.WithValue(nextPageCtx, "action", func(ctx context.Context) context.Context {
						leader := ctx.Value("leader").(string)
						session := ctx.Value("session" + leader).(string)
						name := ctx.Value("name").(string)
						bcclient.LeaderDeselectsMember(session, name)
						return context.WithValue(ctx, "nextPage", nil)
					})
					actionDescCtx := context.WithValue(actionCtx, "actionDesc", "is being deselected?")
					return context.WithValue(actionDescCtx, "currentPage", people)
				}, "is the leader?"),
			},
		},
	}

	list := widgets.NewList()
	return &Emulator{
		list: list,
		pages: map[pageType]page{
			actions: actionPage,
			people:  peoplePage,
		},
		ctx: context.WithValue(context.Background(), "currentPage", actions),
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
