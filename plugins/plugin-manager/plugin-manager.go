package TG

import (
	"fmt"
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

func NewPluginManager() *PluginManagerPlugin {
	return &PluginManagerPlugin{
		plugins: make(map[string]TG.Plugin),
	}
}

func (pm *PluginManagerPlugin) LoadPlugins() {
	pluginDir := "./plugins"
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		log.Fatalf("Error reading plugin directory: %v", err)
	}

	// Iterate over the files in the plugin directory
	for _, file := range files {
		// Only process .so files (plugin files)
		if filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(pluginDir, file.Name())

			// Open the plugin
			plug, err := plugin.Open(pluginPath)
			if err != nil {
				log.Printf("Error loading plugin %s: %v", file.Name(), err)
				continue
			}

			// Look for the "New" function in the plugin
			sym, err := plug.Lookup("New")
			if err != nil {
				log.Printf("Plugin %s does not have a New function: %v", file.Name(), err)
				continue
			}

			// Cast the symbol to a function returning a Plugin
			newPlugin, ok := sym.(func() TG.Plugin)
			if !ok {
				log.Printf("Plugin %s has an invalid New function signature", file.Name())
				continue
			}

			// Initialize the plugin
			pluginInstance := newPlugin()

			// Call the Init function of the plugin
			if apiInstance, ok := pluginInstance.(TG.Plugin); ok {
				// Assuming you have a TG.Api instance
				apiInstance.Init(&Api{})
			} else {
				log.Printf("Plugin %s does not have an Init method", file.Name())
			}

			// Add the plugin to the manager
			pm.AddPlugin(pluginInstance)
			log.Printf("Loaded plugin: %s", file.Name())
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

func (pm *PluginManagerPlugin) RegisterToApi(tg *TG.TG) {
	// This will register the PluginManager into the API
	tg.Api.RegisterCommand("AddPlugin", func(data interface{}) {
		if plugin, ok := data.(TG.Plugin); ok {
			pm.AddPlugin(plugin) // Adjust this based on how plugins are added
		}
	})

	tg.Api.RegisterCommand("GetPlugin", func(data interface{}) {
		if pluginName, ok := data.(string); ok {
			plugin, exists := pm.GetPlugin(pluginName)
			if !exists {
				fmt.Println("Plugin not found.")
			} else {
				fmt.Println("Plugin found:", plugin.Name())
			}
		}
	})
}

func (pm *PluginManagerPlugin) Initialize(tg *TG.TG) {
	// Register the PluginManager with the API
	pm.RegisterToApi(tg)

	// Perform any additional initialization steps if needed
}

func (p *PluginManagerPlugin) OnInstall() {}

func (p *PluginManagerPlugin) OnUninstall() {}
