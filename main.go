package main

import (
	"log"
	"os"
	"plugin"

	TG "github.com/foroughi/tg-edit/tg"
)

func main() {

	// Open or create a log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	tg := TG.NewTG()

	loadPluginManager(tg)

	log.Println("Starting TG-Edit...")

	tg.Event.Dispatch("ON_APP_START", nil)

	tg.Api.Call("Start_UI")

	log.Println("Stopping TG-Edit...")

}

func loadPluginManager(tg *TG.TG) {

	pluginManagerName, exists := tg.Config.Get("pluginmanager")

	if !exists || pluginManagerName == "default" {
		pluginManagerName = "plugin-manager"
	}

	log.Printf("Loading plugin manager: %s...\n", pluginManagerName)
	pluginPath := "./plugins/" + pluginManagerName + ".so"
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Fatalf("Failed to load plugin %s: %v", pluginManagerName, err)
	}

	sym, err := plug.Lookup("New")
	if err != nil {
		log.Fatalf("Plugin %s does not have a New function: %v", pluginManagerName, err)
	}

	newPlugin, ok := sym.(func() TG.Plugin)
	if !ok {
		log.Fatalf("Plugin %s has an invalid New function signature", pluginManagerName)
	}

	pluginInstance := newPlugin()
	pluginInstance.Init(tg)

	log.Printf("Custom plugin %s loaded successfully.", pluginManagerName)

}
