package main

import (
	"log"

	"./internal/api"
	"./internal/plugins"
)

func main() {
	log.Println("Starting TG-Edit...")

	// Create the core API
	api := &api.TGApi{}

	// Load plugins
	plugins.LoadPlugins(api)

	log.Println("TG-Edit is running...")
	select {} // Keep the program running
}
