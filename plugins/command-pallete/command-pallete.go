package main

import (
	"log"

	TG "github.com/foroughi/tg-edit/tg"
)

type CommandPalletePlugin struct {
	tg                     *TG.TG
	commandWindow          any // Pointer to the command palette window
	isCommandPalleteActive bool
	content                string
}

func (p *CommandPalletePlugin) Init(tg *TG.TG) {

	p.tg = tg
	p.isCommandPalleteActive = false
	p.content = ""

	tg.Event.Subscribe("ACTIVE_WINDOW_CHANGED", func(tg *TG.TG, data any) {

		if data == p.commandWindow {

			p.isCommandPalleteActive = true
		} else {
			p.isCommandPalleteActive = false
		}

	})

	tg.Event.Subscribe("ON_KEY", func(tg *TG.TG, data any) {

		if p.isCommandPalleteActive {

			log.Print("====================command pallete key")

			key := data.(string)
			// Handle key input when the command palette is active
			if key == "Esc" {
				// Close the command palette window
				tg.Api.Call("CLOSE_WINDOW", p.commandWindow)
				p.isCommandPalleteActive = false
			} else {
				// Append the key to the command palette content

				p.content += key
				tg.Api.Call("SET_WINDOW_CONTENT", map[string]any{
					"window":  p.commandWindow,
					"content": p.content,
				})
			}
		}
	})
	tg.Api.RegisterCommand("COMMAND", func(tg *TG.TG, data any) {
		p.content = ""
		screenSize := tg.Api.Call("GET_SCREEN_SIZE", nil).(map[string]int)
		screenWidth := screenSize["width"]
		screenHeight := screenSize["height"]
		windowData := map[string]any{
			"title":   "Command Pallete",
			"x":       0,
			"y":       screenHeight - 6, // Position at the bottom of the screen
			"w":       screenWidth,
			"h":       3, // Height of the command palette
			"content": p.content,
		}

		// Save the returned pointer to the command palette window
		p.commandWindow = tg.Api.Call("OPEN_WINDOW", windowData)
		p.isCommandPalleteActive = true
	})

	tg.Key.RegisterKey(":", "COMMAND")
}

func New() TG.Plugin {
	return &CommandPalletePlugin{}
}

func (p *CommandPalletePlugin) Name() string {
	return "CommandPallete"
}

func (p *CommandPalletePlugin) OnInstall() {}

func (p *CommandPalletePlugin) OnUninstall() {}

func (p *CommandPalletePlugin) DependsOn() []string {
	return []string{"UIManager"}
}
