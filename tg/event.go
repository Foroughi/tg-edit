package TG

import (
	"sync"
)

type Subscriber func(event string)

type EventManager struct {
	subscriptions map[string][]Subscriber
	mu            sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscriptions: make(map[string][]Subscriber),
	}
}

func (em *EventManager) Register(name string) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if _, exists := em.subscriptions[name]; !exists {
		em.subscriptions[name] = []Subscriber{}
	}
}

func (em *EventManager) Subscribe(name string, subscriber Subscriber) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, exists := em.subscriptions[name]; !exists {
		em.subscriptions[name] = []Subscriber{}
	}
	em.subscriptions[name] = append(em.subscriptions[name], subscriber)
}

func (em *EventManager) Dispatch(event string) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if subs, exists := em.subscriptions[event]; exists {
		for _, subscriber := range subs {
			subscriber(event)
		}
	}
}
