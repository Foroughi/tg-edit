package main

import (
	TG "github.com/foroughi/tg-edit/tg"
)

type HelloWorldPlugin struct{}

func (p *HelloWorldPlugin) Init(core *TG.TG) {

}

func (p *HelloWorldPlugin) OnInstall() {}

func (p *HelloWorldPlugin) OnUninstall() {}

func (p *HelloWorldPlugin) Name() string {
	return "HelloWorldPlugin"
}

func New() TG.Plugin {
	return &HelloWorldPlugin{}
}
