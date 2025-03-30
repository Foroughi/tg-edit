package TG

import "os"

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

	configManager.Load()
	keyManager.Load(tg)

	for name, action := range defaultCommands {
		eventManager.Register(name, action)
	}

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
	"quit": func(data interface{}) {
		os.Exit(0)
	},
}
