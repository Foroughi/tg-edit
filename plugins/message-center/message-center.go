package main

import (
	"log"

	TG "github.com/foroughi/tg-edit/tg"
)

type Message struct {
	Level        string
	Content      string
	Initiallator string
}

type MessageCenterPlugin struct {
	messages []Message
	tg       *TG.TG
}

func (p *MessageCenterPlugin) Init(tg *TG.TG) {

	p.tg = tg

	p.tg.Api.RegisterCommand("AddMessage", p.AddMessage)
	p.tg.Api.RegisterCommand("GetMessages", p.GetMessages)

	p.tg.Event.Register("MESSAGE_ADDED")

	log.Println("MessageCenter Plugin Initialized")
}

func New() TG.Plugin {
	return &MessageCenterPlugin{
		messages: []Message{},
	}
}

func (p *MessageCenterPlugin) AddMessage(tg *TG.TG, args ...any) any {
	level, content := args[0].(string), args[1].(string)
	msg := Message{Level: level, Content: content}
	p.messages = append(p.messages, msg)

	p.tg.Event.Dispatch("MESSAGE_ADDED", msg)
	log.Printf("[%s] %s", level, content)

	return msg
}

func (p *MessageCenterPlugin) GetMessages(tg *TG.TG, args ...any) any {
	return p.messages
}

func (p *MessageCenterPlugin) Name() string {
	return "MessageCenter"
}

func (p *MessageCenterPlugin) OnInstall() {}

func (p *MessageCenterPlugin) OnUninstall() {}

func (p *MessageCenterPlugin) DependsOn() []string {
	return []string{}
}
