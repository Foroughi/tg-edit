package main

import (
	TG "github.com/foroughi/tg-edit/tg"
	"github.com/gdamore/tcell/v2"
)

type HelloWorldPlugin struct{}

func (p *HelloWorldPlugin) Init(api *TG.Api) {

	x, y := api.GetScreenSize()

	api.WriteText("Hello from HelloWorldPlugin!", x/2, y/2, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
}

func (p *HelloWorldPlugin) Name() string {
	return "HelloWorldPlugin"
}

func New() TG.Plugin {
	return &HelloWorldPlugin{}
}
