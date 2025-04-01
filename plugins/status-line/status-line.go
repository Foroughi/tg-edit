package main

import (
	"log"

	TG "github.com/foroughi/tg-edit/tg"
)

type StatusLinePlugin struct {
	tg               *TG.TG
	statusLineWindow any // Pointer to the command palette window
	leftContent      string
	rightContent     string
	centerContent    string
}

func (p *StatusLinePlugin) Init(tg *TG.TG) {

	p.tg = tg
	p.leftContent = ""
	p.rightContent = ""
	p.centerContent = ""

	tg.Event.Subscribe("ON_UI_START", func(tg *TG.TG, data any) {

		screenSize := tg.Api.Call("GET_SCREEN_SIZE", nil).(map[string]int)
		screenWidth := screenSize["width"]
		screenHeight := screenSize["height"]
		windowData := map[string]any{
			"title":   "Status line",
			"x":       0,
			"y":       screenHeight - 3, // Position at the bottom of the screen
			"w":       screenWidth,
			"h":       3, // Height of the command palette
			"content": p.leftContent + "|" + p.centerContent + "|" + p.rightContent,
		}

		// Save the returned pointer to the command palette window
		p.statusLineWindow = tg.Api.Call("OPEN_WINDOW", windowData)

		tg.Event.Subscribe("ON_KEY", func(tg *TG.TG, data any) {
			log.Print("====================status line key")
			key := data.(string)
			p.rightContent += key
			tg.Api.Call("SET_WINDOW_CONTENT", map[string]any{
				"window":  p.statusLineWindow,
				"content": p.leftContent + "|" + p.centerContent + "|" + p.rightContent,
			})
		})
	})

}

func New() TG.Plugin {
	return &StatusLinePlugin{}
}

func (p *StatusLinePlugin) Name() string {
	return "StatusLine"
}

func (p *StatusLinePlugin) OnInstall() {}

func (p *StatusLinePlugin) OnUninstall() {}

func (p *StatusLinePlugin) DependsOn() []string {
	return []string{"UIManager"}
}
