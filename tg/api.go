package TG

import "sync"

type ApiBridge struct {
	commands map[string]interface{}
	mu       sync.RWMutex
}

func NewApiBridge() *ApiBridge {
	return &ApiBridge{}
}

func (api *ApiBridge) RegisterCommand(name string, fn interface{}) {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.commands[name] = fn
}

func (api *ApiBridge) Call(name string, args ...interface{}) interface{} {
	api.mu.RLock()
	defer api.mu.RUnlock()
	if fn, exists := api.commands[name]; exists {
		if f, ok := fn.(func() string); ok {
			return f()
		} else if f, ok := fn.(func()); ok {
			f()
			return nil
		}
	}
	return nil
}
