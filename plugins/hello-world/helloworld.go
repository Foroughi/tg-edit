package main

import (
	"./internal/api"
	"log"
)

// HelloWorldPlugin is a test plugin that prints a message on startup.
type HelloWorldPlugin struct{}

// Init runs when the plugin is loaded.
func (p *HelloWorldPlugin) Init(api *api.TGApi) {
	api.ShowMessage("Hello from HelloWorldPlugin!")
}

// Name returns the name of the plugin.
func (p *HelloWorldPlugin) Name() string {
	return "HelloWorldPlugin"
}

// New function returns a new instance of the plugin.
func New() Plugin {
	return &HelloWorldPlugin{}
}
