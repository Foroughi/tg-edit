package TG

import (
	"log"

	TG "github.com/foroughi/tg-edit/tg"
)

type MessageLevel string

const (
	Info    MessageLevel = "INFO"
	Warning MessageLevel = "WARNING"
	Error   MessageLevel = "ERROR"
)

type Message struct {
	Level        MessageLevel
	Content      string
	Initiallator string
}

type MessageCenterPlugin struct {
	messages []Message
	tg  *TG.TG
}

func (p *MessageCenterPlugin) Init(tg *TG.TG) {

	p.tg = tg;

	p.tg.Api.RegisterCommand("AddMessage", p.AddMessage)
	p.tg.Api.RegisterCommand("GetMessages", p.GetMessages)

	log.Println("MessageCenter Plugin Initialized")
}

func (p *MessageCenterPlugin) AddMessage(args ...interface{}) interface{} {
	level, content := args[0].(MessageLevel), args[1].(string)
	msg := Message{Level: level, Content: content}
	p.messages = append(p.messages, msg)

	p.tg.Event.Dispatch("message_added", msg)

	return msg
}

func (p *MessageCenterPlugin) GetMessages(args ...interface{}) interface{} {
	return p.messages
}

func (p *MessageCenterPlugin) Name() string {
	return "MessageCenterPlugin"
}

func (p *MessageCenterPlugin) OnInstall() {}

func (p *MessageCenterPlugin) OnUninstall() {}