package TG

import (
	"reflect"
	"sync"
)

type ApiBridge struct {
	commands map[string]interface{}
	mu       sync.RWMutex
}

func NewApiBridge() *ApiBridge {
	return &ApiBridge{
		commands: make(map[string]interface{}),
		mu:       sync.RWMutex{},
	}
}

func (api *ApiBridge) RegisterCommand(name string, fn interface{}) {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.commands[name] = fn
}

func (api *ApiBridge) Call(name string, args ...interface{}) interface{} {
	api.mu.RLock()
	fn, exists := api.commands[name]
	api.mu.RUnlock()
	if !exists {
		return nil
	}

	fnValue := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	out := fnValue.Call(in)

	if len(out) > 0 {
		return out[0].Interface()
	}

	return nil
}
