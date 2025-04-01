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
	args = append([]any{api.tg}, args...) // Ensure the first argument is TG instance

	api.mu.RLock()
	fn, exists := api.commands[name]
	api.mu.RUnlock()
	if !exists {
		log.Printf("[ERROR] Command not found: %s", name)
		return nil
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()
	expectedArgs := fnType.NumIn()

	log.Printf("Calling function: %s", name)
	// Ensure provided arguments match the function's expected parameters
	for len(args) < expectedArgs {
		argType := fnType.In(len(args))
		var zeroValue reflect.Value

		if argType.Kind() == reflect.Ptr {
			zeroValue = reflect.New(argType.Elem()) // Create a pointer to zero value

		} else {
			zeroValue = reflect.Zero(argType) // Standard zero value for the type

		}

		args = append(args, zeroValue.Interface())
	}

	// Convert args to reflect.Value slice
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		if arg == nil {
			argType := fnType.In(i)
			if argType.Kind() == reflect.Ptr {
				in[i] = reflect.New(argType.Elem()) // Create a valid nil pointer

			} else {
				in[i] = reflect.Zero(argType) // Create a zero value

			}
		} else {
			in[i] = reflect.ValueOf(arg)

		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Panic in %s: %v", name, r)
		}
	}()

	// Call function
	out := fnValue.Call(in)
	log.Printf("Function %s executed successfully", name)

	if len(out) > 0 {

		return out[0].Interface()
	}

	return nil
}
