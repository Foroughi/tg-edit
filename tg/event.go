package TG

import (
	"sync"
)

type Event func(data interface{})

type EventManager struct {
	subscriptions map[string]map[int]Event
	lock          sync.RWMutex
	counter       int
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscriptions: make(map[string]map[int]Event),
	}
}

func (em *EventManager) Register(event string, handler Event) {
	em.lock.Lock()
	defer em.lock.Unlock()

	if _, exists := em.subscriptions[event]; !exists {
		em.subscriptions[event] = make(map[int]Event)
	}

	em.counter++
	em.subscriptions[event][em.counter] = handler
}

func (em *EventManager) Dispatch(event string, data interface{}) {
	em.lock.RLock()
	subscriptions, exists := em.subscriptions[event]
	em.lock.RUnlock()
	if !exists {
		return
	}

	for _, subscription := range subscriptions {
		subscription(data)
	}
}
