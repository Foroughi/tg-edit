package main

import (
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
		// Retrieve screen size
		screenSize := tg.Api.Call("GET_SCREEN_SIZE", nil).(map[string]int)
		screenWidth := screenSize["width"]
		screenHeight := screenSize["height"]

		// Dynamically set the style for the status line
		tg.Api.Call("SET_STYLES", map[string]any{
			"key": "status_line",
			"style": map[string]any{
				"fg":        "white",
				"bg":        "black",
				"bold":      true,
				"italic":    false,
				"underline": false,
				"padding":   [4]int{0, 1, 0, 1}, // Top, Right, Bottom, Left
				"margin":    [4]int{0, 0, 0, 0},
				"x":         0,                // Start at the left edge
				"y":         screenHeight - 3, // Position at the bottom of the screen
				"w":         screenWidth,      // Full screen width
				"h":         3,                // Height of the status line
			},
		})

		// Prepare window data with the new style
		windowData := map[string]any{
			"title":   "Status line",
			"content": p.leftContent + "|" + p.centerContent + "|" + p.rightContent,
			"style":   "status_line", // Use the dynamically set style
		}

		// Save the returned pointer to the status line window
		p.statusLineWindow = tg.Api.Call("OPEN_WINDOW", windowData)

		tg.Event.Subscribe("ON_KEY_COMBINATION_FOUND", func(tg *TG.TG, data any) {
			p.rightContent = ""
			p.update()
		})

		tg.Event.Subscribe("ON_KEY_COMBINATION_PROCCESSING", func(tg *TG.TG, data any) {
			key := data.(string)
			p.rightContent = key
			p.update()
		})
	})
}

func (p *StatusLinePlugin) update() {
	p.tg.Api.Call("SET_WINDOW_CONTENT", map[string]any{
		"window":  p.statusLineWindow,
		"content": p.leftContent + "|" + p.centerContent + "|" + p.rightContent,
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
