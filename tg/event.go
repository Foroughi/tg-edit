package TG

import (
	"log"
	"sync"
)

type Event func(tg *TG, data any)

type EventManager struct {
	subscriptions map[string]map[int]Event
	lock          sync.RWMutex
	counter       int
	tg            *TG
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscriptions: make(map[string]map[int]Event),
	}
}

func (em *EventManager) Load(tg *TG) {
	em.tg = tg
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

func (em *EventManager) Dispatch(event string, args any) {

	em.lock.RLock()
	subscriptions, exists := em.subscriptions[event]
	em.lock.RUnlock()
	if !exists {
		return
	}

	for _, subscription := range subscriptions {
		subscription(em.tg, args)
	}
}

func (em *EventManager) Subscribe(event string, handler Event) int {

	log.Printf("%s %d", event, len(em.subscriptions[event]))

	em.lock.Lock()
	defer em.lock.Unlock()

	if _, exists := em.subscriptions[event]; !exists {

		em.subscriptions[event] = make(map[int]Event)
	}

	em.counter++
	em.subscriptions[event][em.counter] = handler
	return em.counter
}
