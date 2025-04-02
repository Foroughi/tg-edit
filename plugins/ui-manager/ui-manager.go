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
	style   string // Field to store the style key
}

type UIManagerPlugin struct {
	screen       tcell.Screen
	windows      []*window
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
	content, _ := windowData["content"].(string)
	style, _ := windowData["style"].(string)

	// Set default style if none is provided
	if style == "" {
		style = "default.win"
	}

	// Retrieve the style to extract default x, y, w, h
	styleData := ui.tg.Api.Call("GET_STYLES", style)
	styleMap, ok := styleData.(map[string]any)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid style format for "+style)
		return nil
	}

	// Extract default values from the style
	defaultX, _ := styleMap["x"].(int)
	defaultY, _ := styleMap["y"].(int)
	defaultW, _ := styleMap["w"].(int)
	defaultH, _ := styleMap["h"].(int)

	// Use provided values or fall back to defaults from the style
	x, _ := windowData["x"].(int)
	if x == 0 {
		x = defaultX
	}

	y, _ := windowData["y"].(int)
	if y == 0 {
		y = defaultY
	}

	w, _ := windowData["w"].(int)
	if w == 0 {
		w = defaultW
	}

	h, _ := windowData["h"].(int)
	if h == 0 {
		h = defaultH
	}

	order := len(ui.windows)

	newWindow := &window{
		title:   title,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		content: content,
		order:   order,
		hidden:  false,
		style:   style, // Store the style key
	}

	ui.windows = append(ui.windows, newWindow)

	ui.tg.Api.Call("AddMessage", "INFO", "Opening "+title)

	ui.activeWindow = newWindow

	ui.draw()

	return newWindow
}

// Function to handle closing a window
func (ui *UIManagerPlugin) closeWindow(data any) any {
	windowPtr, ok := data.(*window)
	if !ok {
		ui.tg.Api.Call("AddMessage", "ERROR", "Invalid data format for CLOSE_WINDOW")
		return nil
	}

	for i := range ui.windows {
		if ui.windows[i] == windowPtr {
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

		if win == windowPtr {

			ui.activeWindow = ui.windows[i]
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

	for _, win := range ui.windows {
		if win == windowPtr { // Compare pointers directly
			win.content = content
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
		if win == windowPtr {
			return win.content
		}
	}

	ui.tg.Api.Call("AddMessage", "ERROR", "Window not found for GET_WINDOW_CONTENT")
	return nil
}

// Refactor RegisterCommand calls to include the new commands
func (ui *UIManagerPlugin) Init(tg *TG.TG) {
	ui.tg = tg
	ui.windows = []*window{}
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

	// Register the new command for styling text
	tg.Api.RegisterCommand("STYLE_TEXT", func(tg *TG.TG, data any) any {
		params, ok := data.(map[string]any)
		if !ok {
			ui.tg.Api.Call("AddMessage", "ERROR", "Invalid data format for STYLE_TEXT")
			return nil
		}

		// Extract the text to style
		text, ok := params["text"].(string)
		if !ok {
			ui.tg.Api.Call("AddMessage", "ERROR", "Invalid text format for STYLE_TEXT")
			return nil
		}

		// Extract the style key
		styleKey, ok := params["style"].(string)
		if !ok {
			ui.tg.Api.Call("AddMessage", "ERROR", "Invalid style key format for STYLE_TEXT")
			return nil
		}

		// Retrieve the style data using GET_STYLES
		styleData := ui.tg.Api.Call("GET_STYLES", styleKey)
		style, ok := styleData.(map[string]any)
		if !ok {
			ui.tg.Api.Call("AddMessage", "ERROR", "Style not found or invalid for key: "+styleKey)
			return text // Return the original text if the style is invalid
		}

		// Use ApplyStyleToText to style the text
		styledText := ui.ApplyStyleToText(text, style)
		return styledText
	})
}

func (ui *UIManagerPlugin) draw() {
	ui.screen.Clear()

	// Sort windows by their order
	orderedWindows := make([]*window, len(ui.windows))
	copy(orderedWindows, ui.windows)

	// Draw all windows except the active one
	for _, win := range orderedWindows {
		if ui.activeWindow != nil && win == ui.activeWindow {
			continue // Skip the active window for now
		}
		if win.hidden {
			continue
		}
		ui.drawWindow(win, win.style) // Pass the window's style
	}

	// Draw the active window last with "default.win.[selected]" as its style
	if ui.activeWindow != nil {
		activeStyle := "default.win.[selected]"
		ui.drawWindow(ui.activeWindow, activeStyle) // Pass the modified style
	}

	ui.screen.Show()
}

// Helper function to draw a single window
func (ui *UIManagerPlugin) drawWindow(win *window, styleKey string) {
	// Retrieve the style using the GET_STYLES command
	styleData := ui.tg.Api.Call("GET_STYLES", styleKey)
	style, ok := styleData.(map[string]any)
	if !ok {
		// If the styleKey does not exist, fall back to "default.win"
		ui.tg.Api.Call("AddMessage", "WARNING", "Style not found: "+styleKey+", falling back to default.win")
		styleData = ui.tg.Api.Call("GET_STYLES", "default.win")
		style, ok = styleData.(map[string]any)
		if !ok {
			ui.tg.Api.Call("AddMessage", "ERROR", "Invalid style format for default.win")
			style = map[string]any{} // Fallback to an empty style
		}
	}

	// Normalize padding and margin
	padding := [4]int{0, 0, 0, 0}
	margin := [4]int{0, 0, 0, 0}
	if p, ok := style["padding"].([]int); ok {
		padding = ui.normalizePaddingMargin(p)
	}
	if m, ok := style["margin"].([]int); ok {
		margin = ui.normalizePaddingMargin(m)
	}

	// Extract other style properties
	fgColor := tcell.ColorWhite
	bgColor := tcell.ColorBlack
	bold := false
	underline := false
	if fg, ok := style["fg"].(string); ok {
		fgColor = tcell.GetColor(fg)
	}
	if bg, ok := style["bg"].(string); ok {
		bgColor = tcell.GetColor(bg)
	}
	if b, ok := style["bold"].(bool); ok {
		bold = b
	}
	if u, ok := style["underline"].(bool); ok {
		underline = u
	}

	// Create the tcell.Style
	tcellStyle := tcell.StyleDefault.Foreground(fgColor).Background(bgColor)
	if bold {
		tcellStyle = tcellStyle.Bold(true)
	}
	if underline {
		tcellStyle = tcellStyle.Underline(true)
	}

	// Adjust window position and size based on margin
	x := win.x + margin[3]             // Left margin
	y := win.y + margin[0]             // Top margin
	w := win.w - margin[1] - margin[3] // Width reduced by left and right margins
	h := win.h - margin[0] - margin[2] // Height reduced by top and bottom margins

	// Check if the border style exists
	borderStyle, hasBorder := style["border"].(map[string]any)
	if hasBorder {
		borderFg := tcell.ColorWhite
		borderBg := tcell.ColorBlack
		borderBold := false
		if fg, ok := borderStyle["fg"].(string); ok {
			borderFg = tcell.GetColor(fg)
		}
		if bg, ok := borderStyle["bg"].(string); ok {
			borderBg = tcell.GetColor(bg)
		}
		if b, ok := borderStyle["bold"].(bool); ok {
			borderBold = b
		}

		borderTcellStyle := tcell.StyleDefault.Foreground(borderFg).Background(borderBg)
		if borderBold {
			borderTcellStyle = borderTcellStyle.Bold(true)
		}

		// Draw window border
		for i := 0; i < w; i++ {
			ui.screen.SetContent(x+i, y, tcell.RuneHLine, nil, borderTcellStyle)
			ui.screen.SetContent(x+i, y+h-1, tcell.RuneHLine, nil, borderTcellStyle)
		}
		for i := 0; i < h; i++ {
			ui.screen.SetContent(x, y+i, tcell.RuneVLine, nil, borderTcellStyle)
			ui.screen.SetContent(x+w-1, y+i, tcell.RuneVLine, nil, borderTcellStyle)
		}
		ui.screen.SetContent(x, y, tcell.RuneULCorner, nil, borderTcellStyle)
		ui.screen.SetContent(x+w-1, y, tcell.RuneURCorner, nil, borderTcellStyle)
		ui.screen.SetContent(x, y+h-1, tcell.RuneLLCorner, nil, borderTcellStyle)
		ui.screen.SetContent(x+w-1, y+h-1, tcell.RuneLRCorner, nil, borderTcellStyle)
	}

	// Adjust content area based on padding
	contentX := x + 1 + padding[3]              // Left padding
	contentY := y + 1 + padding[0]              // Top padding
	contentW := w - 2 - padding[1] - padding[3] // Width reduced by left and right padding
	contentH := h - 2 - padding[0] - padding[2] // Height reduced by top and bottom padding

	// Fill the entire content area with the background style
	for cy := contentY; cy < contentY+contentH; cy++ {
		for cx := contentX; cx < contentX+contentW; cx++ {
			ui.screen.SetContent(cx, cy, ' ', nil, tcellStyle)
		}
	}

	// Check if the title style exists
	titleStyle, hasTitle := style["title"].(map[string]any)
	if hasTitle {
		titleFg := tcell.ColorYellow
		titleBg := tcell.ColorBlack
		titleBold := false
		if fg, ok := titleStyle["fg"].(string); ok {
			titleFg = tcell.GetColor(fg)
		}
		if bg, ok := titleStyle["bg"].(string); ok {
			titleBg = tcell.GetColor(bg)
		}
		if b, ok := titleStyle["bold"].(bool); ok {
			titleBold = b
		}

		titleTcellStyle := tcell.StyleDefault.Foreground(titleFg).Background(titleBg)
		if titleBold {
			titleTcellStyle = titleTcellStyle.Bold(true)
		}

		// Draw window title
		titleRunes := []rune(win.title)
		for i, r := range titleRunes {
			if contentX+i < x+w-1 {
				ui.screen.SetContent(contentX+i, y, r, nil, titleTcellStyle)
			}
		}
	}

	// Draw window content
	contentRunes := []rune(win.content)
	for i, r := range contentRunes {
		cx := contentX + (i % contentW)
		cy := contentY + (i / contentW)
		if cy < contentY+contentH {
			ui.screen.SetContent(cx, cy, r, nil, tcellStyle)
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

// Normalize padding or margin to 4 values
func (ui *UIManagerPlugin) normalizePaddingMargin(values []int) [4]int {
	if len(values) == 1 {
		return [4]int{values[0], values[0], values[0], values[0]} // All sides equal
	} else if len(values) == 2 {
		return [4]int{values[0], values[1], values[0], values[1]} // Top/Bottom, Left/Right
	} else if len(values) == 4 {
		return [4]int{values[0], values[1], values[2], values[3]} // Top, Right, Bottom, Left
	}
	// Default to 0 if invalid
	return [4]int{0, 0, 0, 0}
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
	return []string{"MessageCenter", "HighLight"}
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

func (ui *UIManagerPlugin) ApplyStyleToText(text string, style map[string]any) string {
	// Check if the style is nil or empty, and fall back to "default.text"
	if style == nil || len(style) == 0 {
		styleData := ui.tg.Api.Call("GET_STYLES", "default.text")
		defaultStyle, ok := styleData.(map[string]any)
		if !ok {
			ui.tg.Api.Call("AddMessage", "ERROR", "Invalid style format for default.text")
			defaultStyle = map[string]any{} // Fallback to an empty style
		}
		style = defaultStyle
	}

	// Extract style properties from the map
	fgColor := tcell.ColorWhite
	bgColor := tcell.ColorBlack
	bold := false
	italic := false
	underline := false

	if fg, ok := style["fg"].(string); ok {
		fgColor = tcell.GetColor(fg)
	}
	if bg, ok := style["bg"].(string); ok {
		bgColor = tcell.GetColor(bg)
	}
	if b, ok := style["bold"].(bool); ok {
		bold = b
	}
	if i, ok := style["italic"].(bool); ok {
		italic = i
	}
	if u, ok := style["underline"].(bool); ok {
		underline = u
	}

	// Create a tcell.Style based on the extracted properties
	tcellStyle := tcell.StyleDefault.Foreground(fgColor).Background(bgColor)
	if bold {
		tcellStyle = tcellStyle.Bold(true)
	}
	if italic {
		// Note: tcell does not support italic directly
	}
	if underline {
		tcellStyle = tcellStyle.Underline(true)
	}

	// Apply the style to the text
	styledText := ""
	for _, r := range text {
		// Simulate styled text by appending styled characters
		styledText += string(r) // You can customize this further if needed
	}

	return styledText
}
