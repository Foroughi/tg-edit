package TG

import (
	"log"
)

type KeyManager struct {
	currentSequence  []string          // Tracks the current sequence of keys pressed
	registeredEvents map[string]func() // Registered commands to be executed on key press
}

// Initialize the key manager and set default key combinations
func NewKeyManager() *KeyManager {

	return &KeyManager{
		currentSequence:  []string{},
		registeredEvents: make(map[string]func()),
	}
}

func (km *KeyManager) Load(tg *TG) {
	tg.Event.Subscribe("ON_KEY", km.handleKeyEvent)
}

// Handle a key event: check if it matches a key sequence
func (km *KeyManager) handleKeyEvent(data interface{}) {
	key, ok := data.(string)
	if !ok {
		log.Println("Invalid key data")
		return
	}

	// Add the key to the current sequence
	km.currentSequence = append(km.currentSequence, key)

	// Check if the current sequence matches any of the default key combinations
	if command, exists := km.matchSequence(km.currentSequence); exists {
		// Sequence matched, call the associated command
		if handler, cmdExists := km.registeredEvents[command]; cmdExists {
			handler() // Execute the command
		} else {
			log.Printf("No handler registered for command: %s", command)
		}
		// Clear the sequence after execution
		km.currentSequence = []string{}
	} else {
		// If no match, continue collecting keys for the sequence
		log.Printf("Current sequence: %v", km.currentSequence)
	}
}

// Match the current sequence with the default keys
func (km *KeyManager) matchSequence(seq []string) (string, bool) {
	// Join the sequence to form a key string (e.g., "gx" from ["g", "x"])
	seqStr := ""
	for _, key := range seq {
		seqStr += key
	}

	// Check if the sequence matches any key in the defaultKeys map
	if command, exists := defaultKeys[seqStr]; exists {
		return command, true
	}

	// Check if the sequence is part of a key combination (e.g., "g" as a group)
	if partialCommand, exists := defaultKeys[seq[0]]; exists && partialCommand == "" {
		// The sequence is part of a multi-step command
		return "", false
	}

	return "", false
}
