package main

import (
	"log"

	TG "github.com/foroughi/tg-edit/tg"
	"github.com/gdamore/tcell/v2"
)

type UIManagerPlugin struct {
	screen      tcell.Screen
	tg          *TG.TG
	exitChannel chan struct{}
}

func (ui *UIManagerPlugin) Init(tg *TG.TG) {
	ui.tg = tg
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	ui.screen = screen

	tg.Event.Subscribe("ON_Quit", func(tg *TG.TG, args any) {
		log.Print("calling ***************")
		ui.exitChannel <- struct{}{}
	})

	tg.Api.RegisterCommand("Start_UI", func(tg *TG.TG, data any) {
		defer ui.screen.Fini()
		if err := ui.screen.Init(); err != nil {
			log.Fatalf("Failed to initialize screen: %v", err)
		}
		ui.eventLoop()
	})

}

func (ui *UIManagerPlugin) eventLoop() {
	ui.screen.Clear()
	ui.exitChannel = make(chan struct{})
	go func() {
		for {
			ev := ui.screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:

				if ev.Key() == tcell.KeyEscape {

					ui.exitChannel <- struct{}{}
					return
				}

				ui.tg.Event.Dispatch("ON_KEY", ev)
			}
			ui.screen.Show()
		}
	}()
	<-ui.exitChannel
}

func (p *UIManagerPlugin) Name() string {
	return "UIManager"
}

func New() TG.Plugin {
	return &UIManagerPlugin{}
}

func (p *UIManagerPlugin) OnInstall() {}

func (p *UIManagerPlugin) OnUninstall() {}

func (p *UIManagerPlugin) DependsOn() []string {
	return []string{"MessageCenter"}
}

func (ui *UIManagerPlugin) DrawText(args ...any) any {

	if len(args) < 3 {

		return ui.tg.Api.Call("AddMessage", "ERROR", "Usage: ui.draw_text <x> <y> <text>")

	}

	x, y, text := args[0].(int), args[1].(int), args[2].(string)
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	runes := []rune(text)
	for i, r := range runes {
		ui.screen.SetContent(x+i, y, r, nil, style)
	}

	return nil
}
