package main

import (
	"strings"

	TG "github.com/foroughi/tg-edit/tg"
)

// HighLightPlugin with decoupled styles
type HighLightPlugin struct {
	tg     *TG.TG
	styles map[string]any // Recursive map for unlimited levels
}

func (p *HighLightPlugin) Init(tg *TG.TG) {
	p.tg = tg
	p.styles = map[string]any{
		"default": map[string]any{
			"text": map[string]any{
				"fg":        "white",
				"bg":        "black",
				"bold":      false,
				"italic":    false,
				"underline": false,
				"padding":   [4]int{0},
				"margin":    [4]int{0},
			},
			"win": map[string]any{
				"w": 80,
				"h": 25,
				"x": 0,
				"y": 0,
				"border": map[string]any{
					"fg":   "white",
					"bg":   "black",
					"bold": true,
				},
				"title": map[string]any{
					"fg":   "white",
					"bg":   "black",
					"bold": true},
				"padding": [4]int{0},
				"margin":  [4]int{0},
				"[selected]": map[string]any{
					"Width":  80,
					"Height": 25,
					"X":      0,
					"Y":      0,
					"border": map[string]any{
						"fg":   "white",
						"bg":   "black",
						"Bold": true,
					},
					"title": map[string]any{
						"fg":   "white",
						"bg":   "black",
						"bold": true},
					"padding": [4]int{0},
					"margin":  [4]int{0},
				},
			},
		},
	}

	// Register the Get_STYLES command
	tg.Api.RegisterCommand("GET_STYLES", func(tg *TG.TG, data any) any {
		keys, ok := data.(string)
		if !ok {
			tg.Api.Call("AddMessage", "ERROR", "Invalid data format for Get_STYLES")
			return nil
		}
		return p.getStyle(splitKeyPath(keys))
	})

	// Register the Set_STYLES command
	tg.Api.RegisterCommand("SET_STYLES", func(tg *TG.TG, data any) any {
		params, ok := data.(map[string]any)
		if !ok {
			tg.Api.Call("AddMessage", "ERROR", "Invalid data format for Set_STYLES")
			return nil
		}

		key, ok := params["key"].(string)
		if !ok {
			tg.Api.Call("AddMessage", "ERROR", "Invalid keys for Set_STYLES")
			return nil
		}

		style, ok := params["style"].(map[string]any)
		if !ok {
			tg.Api.Call("AddMessage", "ERROR", "Invalid style for Set_STYLES")
			return nil
		}

		p.setStyle(key, style)

		tg.Api.Call("AddMessage", "INFO", "Style updated: "+key)
		return nil
	})
}

// Recursive function to get a style
func (p *HighLightPlugin) getStyle(keys []string) any {
	current := p.styles
	for i, key := range keys {
		if next, ok := current[key]; ok {
			if i == len(keys)-1 {
				// Last key, should be a style
				return next
			}

			// Not the last key, should be a map
			if nextMap, ok := next.(map[string]any); ok {
				current = nextMap
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
}

// Function to set a style using a dot-separated string
func (p *HighLightPlugin) setStyle(keyPath string, style any) {
	// Split the keyPath into parts
	keys := splitKeyPath(keyPath)

	// Traverse the styles map
	current := p.styles
	for i, key := range keys {
		if i == len(keys)-1 {
			// Last key, set the style
			current[key] = style
			return
		}

		// Not the last key, ensure it's a map
		if next, ok := current[key]; ok {
			if nextMap, ok := next.(map[string]any); ok {
				current = nextMap
			} else {
				return
			}
		} else {
			// Create a new map if the key doesn't exist
			newMap := map[string]any{}
			current[key] = newMap
			current = newMap
		}
	}

}

// Helper function to split a dot-separated key path
func splitKeyPath(keyPath string) []string {
	return strings.Split(keyPath, ".")
}

func New() TG.Plugin {
	return &HighLightPlugin{}
}

func (p *HighLightPlugin) Name() string {
	return "HighLight"
}

func (p *HighLightPlugin) OnInstall() {}

func (p *HighLightPlugin) OnUninstall() {}

func (p *HighLightPlugin) DependsOn() []string {
	return []string{}
}
