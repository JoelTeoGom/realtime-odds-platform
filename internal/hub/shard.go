package hub

import (
	"sync"

	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/ports/outbound"
)

type Shard struct {
	mu   sync.RWMutex
	subs map[string]map[string]outbound.Client
}

func newShard() *Shard {
	return &Shard{
		subs: make(map[string]map[string]outbound.Client),
		mu:   sync.RWMutex{},
	}
}

func (s *Shard) subscribe(c outbound.Client, topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients, ok := s.subs[topic]
	if !ok {
		clients = make(map[string]outbound.Client)
		s.subs[topic] = clients
	}
	clients[c.ID()] = c
}

func (s *Shard) unsubscribe(clientID, topic string) {
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
