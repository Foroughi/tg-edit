package main

import (
	"log"
	"strings"

	TG "github.com/foroughi/tg-edit/tg"
	"github.com/gdamore/tcell/v2"
)

type window struct {
	title   string
	x       int
	y       int
	w       int
	h       int
	content string
	order   int
	hidden  bool
}

type UIManagerPlugin struct {
	screen       tcell.Screen
	windows      []window
	activeWindow *window
	tg           *TG.TG
	exitFlag     bool
}

// Function to handle opening a window
func (ui *UIManagerPlugin) openWindow(data any) any {
	windowData, ok := data.(map[string]any)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid data format for OPEN_WINDOW")
		return nil
	}

	title, _ := windowData["title"].(string)
	x, _ := windowData["x"].(int)
	y, _ := windowData["y"].(int)
	w, _ := windowData["w"].(int)
	h, _ := windowData["h"].(int)
	content, _ := windowData["content"].(string)
	order := len(ui.windows)

	newWindow := window{
		title:   title,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		content: content,
		order:   order,
		hidden:  false,
	}

	ui.windows = append(ui.windows, newWindow)

	ui.tg.Api.Call("AddMessage", "INFO", "Opening "+title)

	ui.activeWindow = &ui.windows[len(ui.windows)-1]

	ui.draw()

	return &ui.windows[len(ui.windows)-1]
}

// Function to handle closing a window
func (ui *UIManagerPlugin) closeWindow(data any) any {
	windowPtr, ok := data.(*window)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid data format for CLOSE_WINDOW")
		return nil
	}

	for i := range ui.windows {
		if &ui.windows[i] == windowPtr {
			ui.windows = append(ui.windows[:i], ui.windows[i+1:]...)
			if ui.activeWindow == windowPtr {
				ui.activeWindow = nil
			}
			ui.tg.Api.Call("AddMessage", "INFO", "Window closed")
			ui.draw()
			return nil
		}
	}

	ui.tg.Api.Call("AddMessage", "ERROR", "Window not found")
	return nil
}

// Function to handle making a window active
func (ui *UIManagerPlugin) makeWindowActive(data any) any {

	windowPtr, ok := data.(*window)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid data format for MAKE_WINDOW_ACTIVE")
		return nil
	}

	for i, win := range ui.windows {

		if &win == windowPtr {

			ui.activeWindow = &ui.windows[i]
			ui.tg.Event.Dispatch("ACTIVE_WINDOW_CHANGED", data)
			ui.tg.Api.Call("AddMessage", "INFO", "Window set as active")
			ui.draw()
			return nil
		}
	}

	ui.tg.Api.Call("AddMessage", "ERROR", "Window not found to active")
	return nil
}

// Function to get the screen size
func (ui *UIManagerPlugin) getScreenSize(data any) any {
	width, height := ui.screen.Size()
	return map[string]int{"width": width, "height": height}
}

// Function to set the content of a window
func (ui *UIManagerPlugin) setWindowContent(data any) any {
	params, ok := data.(map[string]any)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid data format for SET_WINDOW_CONTENT")
		return nil
	}

	windowPtr, ok := params["window"].(*window)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid window pointer for SET_WINDOW_CONTENT")
		return nil
	}

	content, ok := params["content"].(string)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid content format for SET_WINDOW_CONTENT")
		return nil
	}

	for i := range ui.windows {
		if &ui.windows[i] == windowPtr {
			ui.windows[i].content = content
			ui.draw()
			ui.tg.Api.Call("AddMessage", "INFO", "Window content updated")
			return nil
		}
	}

	ui.tg.Api.Call("AddMessage", "ERROR", "Window not found for SET_WINDOW_CONTENT")
	return nil
}

// Function to get the content of a window
func (ui *UIManagerPlugin) getWindowContent(data any) any {
	windowPtr, ok := data.(*window)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid window pointer for GET_WINDOW_CONTENT")
		return nil
	}

	for _, win := range ui.windows {
		if &win == windowPtr {
			return win.content
		}
	}

	ui.tg.Api.Call("AddMessage", "ERROR", "Window not found for GET_WINDOW_CONTENT")
	return nil
}

// Refactor RegisterCommand calls to include the new commands
func (ui *UIManagerPlugin) Init(tg *TG.TG) {
	ui.tg = tg
	ui.windows = []window{}
	ui.activeWindow = nil

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	ui.screen = screen
	ui.exitFlag = false

	tg.Event.Subscribe("ON_Quit", func(tg *TG.TG, args any) {
		ui.exitFlag = true
	})

	tg.Api.RegisterCommand("Start_UI", func(tg *TG.TG, data any) {
		defer ui.screen.Fini()
		if err := ui.screen.Init(); err != nil {
			log.Fatalf("Failed to initialize screen: %v", err)
		}

		tg.Event.Dispatch("ON_UI_START", nil)

		ui.eventLoop()
	})

	tg.Api.RegisterCommand("SET_WINDOW_CONTENT", func(tg *TG.TG, data any) any {
		return ui.setWindowContent(data)
	})

	tg.Api.RegisterCommand("GET_WINDOW_CONTENT", func(tg *TG.TG, data any) any {
		return ui.getWindowContent(data)
	})

	tg.Api.RegisterCommand("OPEN_WINDOW", func(tg *TG.TG, data any) any {
		return ui.openWindow(data)
	})

	tg.Api.RegisterCommand("CLOSE_WINDOW", func(tg *TG.TG, data any) any {
		return ui.closeWindow(data)
	})

	tg.Api.RegisterCommand("ACTIVE_WINDOW", func(tg *TG.TG, data any) any {
		return ui.makeWindowActive(data)
	})

	tg.Api.RegisterCommand("GET_SCREEN_SIZE", func(tg *TG.TG, data any) any {
		return ui.getScreenSize(data)
	})
}

func (ui *UIManagerPlugin) draw() {
	ui.screen.Clear()

	// Sort windows by their order
	orderedWindows := make([]window, len(ui.windows))
	copy(orderedWindows, ui.windows)

	// Draw all windows except the active one
	for _, win := range orderedWindows {
		if ui.activeWindow != nil && &win == ui.activeWindow {
			continue // Skip the active window for now
		}
		if win.hidden {
			continue
		}
		ui.drawWindow(win)
	}

	// Draw the active window last
	if ui.activeWindow != nil {
		ui.drawWindow(*ui.activeWindow)
	}

	ui.screen.Show()
}

// Helper function to draw a single window
func (ui *UIManagerPlugin) drawWindow(win window) {
	// Draw window border
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	for i := 0; i < win.w; i++ {
		ui.screen.SetContent(win.x+i, win.y, tcell.RuneHLine, nil, style)
		ui.screen.SetContent(win.x+i, win.y+win.h-1, tcell.RuneHLine, nil, style)
	}
	for i := 0; i < win.h; i++ {
		ui.screen.SetContent(win.x, win.y+i, tcell.RuneVLine, nil, style)
		ui.screen.SetContent(win.x+win.w-1, win.y+i, tcell.RuneVLine, nil, style)
	}
	ui.screen.SetContent(win.x, win.y, tcell.RuneULCorner, nil, style)
	ui.screen.SetContent(win.x+win.w-1, win.y, tcell.RuneURCorner, nil, style)
	ui.screen.SetContent(win.x, win.y+win.h-1, tcell.RuneLLCorner, nil, style)
	ui.screen.SetContent(win.x+win.w-1, win.y+win.h-1, tcell.RuneLRCorner, nil, style)

	// Draw window title
	titleRunes := []rune(win.title)
	for i, r := range titleRunes {
		if win.x+1+i < win.x+win.w-1 {
			ui.screen.SetContent(win.x+1+i, win.y, r, nil, style)
		}
	}

	// Draw window content
	contentRunes := []rune(win.content)
	for i, r := range contentRunes {
		x := win.x + 1 + (i % (win.w - 2))
		y := win.y + 1 + (i / (win.w - 2))
		if y < win.y+win.h-1 {
			ui.screen.SetContent(x, y, r, nil, style)
		}
	}
}

func (ui *UIManagerPlugin) eventLoop() {
	ui.screen.Clear()

	for {

		if ui.exitFlag {
			return
		}

		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			ui.tg.Event.Dispatch("ON_KEY", ui.getKeyString(ev))
		}

		ui.draw()

	}
}

func (ui *UIManagerPlugin) getKeyString(ev *tcell.EventKey) string {
	if ev.Modifiers()&tcell.ModCtrl != 0 {
		return "Ctrl+" + ev.Name()
	}
	if ev.Modifiers()&tcell.ModAlt != 0 {
		return "Alt+" + ev.Name()
	}
	return strings.TrimPrefix(strings.TrimSuffix(ev.Name(), "]"), "Rune[") // Default key name
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
