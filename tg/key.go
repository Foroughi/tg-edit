package TG

import (
	"log"
	"strings"
)

type KeyManager struct {
	currentSequence []string // Tracks the current sequence of keys pressed
	tg              *TG
}

// Initialize the key manager and set default key combinations
func NewKeyManager() *KeyManager {

	return &KeyManager{
		currentSequence: []string{},
	}
}

func (km *KeyManager) Load(tg *TG) {
	km.tg = tg
	km.tg.Event.Subscribe("ON_KEY", func(tg *TG, data any) {

		km.handleKeyEvent(data)
	})
}

// Handle a key event: check if it matches a key sequence
func (km *KeyManager) handleKeyEvent(data any) {

	key, ok := data.(string)

	if !ok {

		return
	}

	// Add the key to the current sequence
	km.currentSequence = append(km.currentSequence, key)
	log.Printf("Current key sequence: %v", km.currentSequence)
	// Check if the current sequence matches any of the default key combinations
	if command, exists := km.matchSequence(); exists {
		// Sequence matched, call the associated command

		km.tg.Api.Call(command)

		// Clear the sequence after execution
		km.currentSequence = []string{}
	} else {
		// Clear the sequence after execution
		matched := false
		for key := range defaultKeys {
			if strings.HasPrefix(key, strings.Join(km.currentSequence, "")) {
				matched = true
				break
			}
		}
		if !matched {
			km.currentSequence = []string{}
		}
	}

}

func (km *KeyManager) matchSequence() (string, bool) {
	// Convert currentSequence (e.g., ["Rune[g]", "Rune[x]"]) to a simple string (e.g., "gx")
	var cleanedSeq []string
	for _, key := range km.currentSequence {
		cleanedSeq = append(cleanedSeq, key)
	}
	seqStr := strings.Join(cleanedSeq, "")

	// Check for a direct match
	if command, exists := defaultKeys[seqStr]; exists && command != "" {
		return command, true
	}

	// If it's a group key (e.g., "g" with ""), continue waiting
	if _, exists := defaultKeys[seqStr]; exists {
		return "", false
	}

	return "", false
}

// RegisterKeyCombination allows plugins to register a key combination with a command
func (km *KeyManager) RegisterKey(keyCombination string, command string) {

	defaultKeys[keyCombination] = command
}
