package hub

import (
	"sync"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/ports/outbound"
)

type EventShard struct {
	mu           sync.RWMutex
	eventSubsMap map[string][]string
}

type UserShard struct {
	mu              sync.RWMutex
	clientSocketMap map[string]outbound.Connection
}

func newUserShard() *UserShard {
	return &UserShard{
		clientSocketMap: make(map[string]outbound.Connection),
		mu:              sync.RWMutex{},
	}
}

func newEventShard() *EventShard {
	return &EventShard{
		eventSubsMap: make(map[string][]string),
		mu:           sync.RWMutex{},
	}
}

func (s *UserShard) subscribe(c outbound.Connection, topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients, ok := s.clientSocketMap[topic]
	if !ok {
		clients = make(map[string]outbound.Connection)
		s.subs[topic] = clients
	}
	clients[c.ID()] = c
}

func (s *UserShard) unsubscribe(clientID, topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients, ok := s.subs[topic]
	if !ok {
		return
	}

	delete(clients, clientID)
	if len(clients) == 0 {
		delete(s.subs, topic) // an empty inner map would leak
	}
}
