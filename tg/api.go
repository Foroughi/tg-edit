package TG

import (
	"log"
	"reflect"
	"sync"
)

type ApiBridge struct {
	commands map[string]any
	mu       sync.RWMutex
	tg       *TG
}

func NewApiBridge() *ApiBridge {
	return &ApiBridge{
		commands: make(map[string]any),
		mu:       sync.RWMutex{},
	}
}

func (api *ApiBridge) Load(tg *TG) {
	api.tg = tg
}

func (api *ApiBridge) RegisterCommand(name string, fn any) {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.commands[name] = fn

	log.Print("registering " + name)
}

func (api *ApiBridge) Call(name string, args ...any) any {

	args = append([]any{api.tg}, args...)

	api.mu.RLock()
	fn, exists := api.commands[name]
	api.mu.RUnlock()
	if !exists {
		log.Printf("Command not found: %s", name)
		return nil
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// Ensure at least the required number of arguments are passed
	expectedArgs := fnType.NumIn()
	if len(args) < expectedArgs {
		// Fill missing arguments with `nil`
		for len(args) < expectedArgs {
			args = append(args, nil)
		}
	}

	in := make([]reflect.Value, len(args))

	log.Printf("Calling function: %s with args: %v", name, args)

	// Recover from panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in %s: %v", name, r)
		}
	}()

	out := fnValue.Call(in)
	log.Print("Function call completed")

	if len(out) > 0 {
		return out[0].Interface()
	}

	return nil
}
