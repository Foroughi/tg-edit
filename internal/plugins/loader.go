package plugins

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"

	"./internal/api"
)

// Plugin is the standard interface that all plugins must implement.
type Plugin interface {
	Init(api *api.TGApi)
	Name() string
}

// LoadPlugins dynamically loads all Go plugins from the "plugins" directory.
func LoadPlugins(api *api.TGApi) {
	pluginsDir := "plugins"

	// Check if the plugins directory exists
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		log.Println("No plugins directory found. Skipping plugin loading.")
		return
	}

	// Iterate over all plugin files
	err := filepath.Walk(pluginsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".so" { // Go plugins are compiled as .so files
			log.Println("Loading plugin:", path)
			loadPlugin(path, api)
		}
		return nil
	})

	if err != nil {
		log.Println("Error loading plugins:", err)
	}
}

// loadPlugin loads a single Go plugin
func loadPlugin(path string, api *api.TGApi) {
	p, err := plugin.Open(path)
	if err != nil {
		log.Println("Failed to open plugin:", err)
		return
	}

	// Look for a "New" function that returns a Plugin instance
	sym, err := p.Lookup("New")
	if err != nil {
		log.Println("Plugin does not have a New() function:", err)
		return
	}

	newFunc, ok := sym.(func() Plugin)
	if !ok {
		log.Println("Invalid plugin format")
		return
	}

	instance := newFunc()
	instance.Init(api)
	log.Println("Loaded plugin:", instance.Name())
}
