package main

import (
	"log"
	"os"
	"path/filepath"
	"plugin"

	TG "github.com/foroughi/tg-edit/tg"
)

type PluginManagerPlugin struct {
	plugins map[string]TG.Plugin
	tg      *TG.TG
}

func New() TG.Plugin {
	return &PluginManagerPlugin{
		plugins: make(map[string]TG.Plugin),
	}
}

func (pm *PluginManagerPlugin) shouldLoadPlugin(fileName string) bool {
	// Prevent loading itself
	return fileName != "plugin-manager.so"
}

func (pm *PluginManagerPlugin) LoadPlugins() {
	pluginDir := "./plugins"
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		log.Fatalf("Error reading plugin directory: %v", err)
	}

	// Step 1: Collect all plugins
	pendingPlugins := make(map[string]TG.Plugin)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" && pm.shouldLoadPlugin(file.Name()) {
			pluginPath := filepath.Join(pluginDir, file.Name())

			plug, err := plugin.Open(pluginPath)
			if err != nil {
				log.Printf("Error loading plugin %s: %v", file.Name(), err)
				continue
			}

			sym, err := plug.Lookup("New")
			if err != nil {
				log.Printf("Plugin %s does not have a New function: %v", file.Name(), err)
				continue
			}

			newPlugin, ok := sym.(func() TG.Plugin)
			if !ok {
				log.Printf("Plugin %s has an invalid New function signature", file.Name())
				continue
			}

			pluginInstance := newPlugin()
			pendingPlugins[pluginInstance.Name()] = pluginInstance
		}
	}

	// Step 2: Load plugins in correct order
	loadedPlugins := make(map[string]TG.Plugin)

	for len(pendingPlugins) > 0 {
		progress := false

		for name, plugin := range pendingPlugins {
			missingDeps := []string{}
			for _, dep := range plugin.DependsOn() {
				if _, exists := pendingPlugins[dep]; !exists && loadedPlugins[dep] == nil {
					missingDeps = append(missingDeps, dep)
				}
			}

			if len(missingDeps) > 0 {
				log.Printf("Skipping plugin %s: missing dependencies: %v", name, missingDeps)
				delete(pendingPlugins, name) // Remove it since dependencies are missing
				continue
			}

			// Load the plugin since all dependencies exist
			plugin.Init(pm.tg)
			loadedPlugins[name] = plugin
			pm.AddPlugin(plugin)
			delete(pendingPlugins, name)
			log.Printf("Loaded plugin: %s", name)
			progress = true
		}

		// Step 3: If no progress, circular dependency detected
		if !progress {
			log.Fatalf("Circular dependency detected! Unresolved plugins: %v", pendingPlugins)
		}
	}
}

func (pm *PluginManagerPlugin) AddPlugin(plugin TG.Plugin) {
	if _, exists := pm.plugins[plugin.Name()]; exists {
		log.Printf("Plugin '%s' already exists.", plugin.Name())
		return
	}
	pm.plugins[plugin.Name()] = plugin
	log.Printf("Plugin '%s' added successfully.", plugin.Name())
}

func (pm *PluginManagerPlugin) GetPlugin(name string) (TG.Plugin, bool) {
	plugin, exists := pm.plugins[name]
	return plugin, exists
}

func (pm *PluginManagerPlugin) Init(tg *TG.TG) {

	pm.tg = tg

	pm.tg.Event.Subscribe("ON_APP_START", func(data interface{}) {
		pm.LoadPlugins()

	})

}

func (p *PluginManagerPlugin) Name() string {
	return "PluginManager"
}

func (p *PluginManagerPlugin) OnInstall() {}

func (p *PluginManagerPlugin) OnUninstall() {}

func (p *PluginManagerPlugin) DependsOn() []string {
	return []string{}
}
