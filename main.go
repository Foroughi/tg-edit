package main

import (
	"log"
	"os"
	"path/filepath"

	"plugin"

	TG "github.com/foroughi/tg-edit/tg"
	"github.com/gdamore/tcell/v2"
)

func LoadPlugins(apiInstance *TG.Api) {
	pluginsDir := "plugins"

	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		log.Println("No plugins directory found. Skipping plugin loading.")
		return
	}

	err := filepath.Walk(pluginsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".so" {
			log.Println("Loading plugin:", path)
			loadPlugin(path, apiInstance)
		}
		return nil
	})

	if err != nil {
		log.Println("Error loading plugins:", err)
	}
}

func loadPlugin(path string, apiInstance *TG.Api) {
	p, err := plugin.Open(path)
	if err != nil {
		log.Println("Failed to open plugin:", err)
		return
	}

	sym, err := p.Lookup("New")
	if err != nil {
		log.Println("Plugin does not have a New() function:", err)
		return
	}

	newFunc, ok := sym.(func() TG.Plugin)
	if !ok {
		log.Println("Invalid plugin format")
		return
	}

	instance := newFunc()
	instance.Init(apiInstance)
	log.Println("Loaded plugin:", instance.Name())
}

func main() {

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	screen.Init()
	defer screen.Fini()

	screen.Clear()

	tgApi := &TG.Api{}
	tgApi.Init(screen)

	log.Println("Starting TG-Edit...")

	LoadPlugins(tgApi)

	log.Println("TG-Edit is running...")

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return
			}
		}
	}
}
