package TG

import (
	"strings"
)

type KeyManager struct {
	currentSequence []string // Tracks the current sequence of keys pressed
	tg              *TG
	recording       bool
}

// Initialize the key manager and set default key combinations
func NewKeyManager() *KeyManager {

	return &KeyManager{
		currentSequence: []string{},
	}
}

func (km *KeyManager) Load(tg *TG) {
	km.tg = tg
	km.recording = true

	km.tg.Event.Register("ON_KEY_COMBINATION_FOUND")
	km.tg.Event.Register("ON_KEY_COMBINATION_PROCCESSING")

	tg.Api.RegisterCommand("RECORD_KEYS", func(tg *TG, data any) {
		km.recording = true
	})

	tg.Api.RegisterCommand("DONT_RECORD_KEYS", func(tg *TG, data any) {
		km.recording = false
	})

	km.tg.Event.Subscribe("ON_KEY", func(tg *TG, data any) {
		if km.recording {
			km.handleKeyEvent(data)
		}
	})
}

// Handle a key event: check if it matches a key sequence
func (km *KeyManager) handleKeyEvent(data any) {

	key, ok := data.(string)

	if !ok {
		return
	}

	km.currentSequence = append(km.currentSequence, key)
	if command, exists := km.matchSequence(); exists {

		km.tg.Event.Dispatch("ON_KEY_COMBINATION_FOUND", strings.Join(km.currentSequence, " "))

		km.tg.Api.Call(command)

		km.currentSequence = []string{}
	} else {

		matched := false
		for key := range defaultKeys {
			if strings.HasPrefix(key, strings.Join(km.currentSequence, "")) {
				matched = true
				km.tg.Event.Dispatch("ON_KEY_COMBINATION_PROCCESSING", strings.Join(km.currentSequence, " "))
				break
			}
		}
		if !matched {

			km.tg.Event.Dispatch("ON_KEY_COMBINATION_FOUND", nil)

			km.currentSequence = []string{}
		}
	}

}

func (km *KeyManager) matchSequence() (string, bool) {

	var cleanedSeq []string

	cleanedSeq = append(cleanedSeq, km.currentSequence...)

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
