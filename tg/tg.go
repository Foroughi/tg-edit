package TG

type TG struct {
	Api    *ApiBridge
	Event  *EventManager
	Config *ConfigManager
}

func NewTG() *TG {
	eventManager := NewEventManager()
	apiBridge := NewApiBridge()
	configManager := NewConfigManager()

	configManager.Load()

	tg := &TG{
		Event:  eventManager,
		Api:    apiBridge,
		Config: configManager,
	}

	return tg
}
