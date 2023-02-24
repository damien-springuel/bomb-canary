package main

import (
	"log"

	"github.com/damien-springuel/bomb-canary/cli/emulator"
	"github.com/gizak/termui/v3"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	grid := termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	emulator := emulator.New()
	grid.Set(termui.NewRow(1, emulator))

	termui.Render(grid)

	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(termui.Resize)
			grid.SetRect(0, 0, payload.Width, payload.Height)
			termui.Clear()
		}

		emulator.HandleUiEvent(e)
		termui.Render(grid)
	}
}
