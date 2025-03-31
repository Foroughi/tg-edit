package TG

import (
	"strings"

	"github.com/gdamore/tcell/v2"
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

	keyEvent, ok := data.(*tcell.EventKey)

	if !ok {

		return
	}

	key := km.getKeyString(keyEvent)

	// Add the key to the current sequence
	km.currentSequence = append(km.currentSequence, key)

	// Check if the current sequence matches any of the default key combinations
	if command, exists := km.matchSequence(); exists {
		// Sequence matched, call the associated command

		km.tg.Api.Call(command)

		// Clear the sequence after execution
		km.currentSequence = []string{}
	}
}

func (km *KeyManager) getKeyString(ev *tcell.EventKey) string {
	if ev.Modifiers()&tcell.ModCtrl != 0 {
		return "Ctrl+" + ev.Name()
	}
	if ev.Modifiers()&tcell.ModAlt != 0 {
		return "Alt+" + ev.Name()
	}
	return ev.Name() // Default key name
}

func (km *KeyManager) matchSequence() (string, bool) {
	// Convert currentSequence (e.g., ["Rune[g]", "Rune[x]"]) to a simple string (e.g., "gx")
	var cleanedSeq []string
	for _, key := range km.currentSequence {
		cleanedSeq = append(cleanedSeq, strings.TrimPrefix(strings.TrimSuffix(key, "]"), "Rune["))
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
