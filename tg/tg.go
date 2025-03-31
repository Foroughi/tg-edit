package TG

import "log"

type TG struct {
	Api    *ApiBridge
	Event  *EventManager
	Config *ConfigManager
	Key    *KeyManager
}

func NewTG() *TG {

	eventManager := NewEventManager()
	apiBridge := NewApiBridge()
	configManager := NewConfigManager()
	keyManager := NewKeyManager()

	tg := &TG{
		Event:  eventManager,
		Api:    apiBridge,
		Config: configManager,
		Key:    keyManager,
	}

	for name, action := range defaultCommands {
		tg.Api.RegisterCommand(name, action)
	}

	configManager.Load()
	keyManager.Load(tg)
	apiBridge.Load(tg)
	eventManager.Load(tg)

	return tg
}

// Default configurations
var defaultConfig = map[string]string{
	"pluginmanager": "default",
}

var defaultKeys = map[string]string{
	"g":  "",
	"gx": "quit",
	":":  "openCommandPallet",
}

var defaultCommands = map[string]Event{
	"quit": func(tg *TG, data any) {
		log.Print("efwefwefewfwef")
		tg.Event.Dispatch("ON_Quit", data)
	},
}
