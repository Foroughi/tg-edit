package TG

import (
	"github.com/gdamore/tcell/v2"
)

type Api struct {
	screen tcell.Screen
}

func (api *Api) Init(screen tcell.Screen) {
	api.screen = screen
}

func (api *Api) RunCommand(cmd string) {
}

func (api *Api) GetScreenSize() (int, int) {
	return api.screen.Size()
}

func (api *Api) WriteText(msg string, x int, y int, style tcell.Style) {

	for i, ch := range msg {
		api.screen.SetContent(x+i, y, ch, nil, style)
	}
	api.screen.Show()

}

func (api *Api) ClearScreen() {
	api.screen.Clear()
	api.screen.Show()
}
