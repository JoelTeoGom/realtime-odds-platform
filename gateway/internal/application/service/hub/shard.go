package hub

import (
	"sync"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/ports/outbound"
)

// EventShard holds, per topic, the clients listening to it, keyed by client id
// so subscribe and unsubscribe stay O(1) instead of scanning the topic.
type EventShard struct {
	mu           sync.RWMutex
	eventSubsMap map[string]map[string]outbound.Connection
}

func newEventShard() *EventShard {
	return &EventShard{
		eventSubsMap: make(map[string]map[string]outbound.Connection),
	}
}

// subscribe adds the client to the topic, creating the topic on first use.
// Re-subscribing with the same id is a no-op that refreshes the connection.
func (s *EventShard) subscribe(c outbound.Connection, topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients, ok := s.eventSubsMap[topic]
	if !ok {
		clients = make(map[string]outbound.Connection)
		s.eventSubsMap[topic] = clients
	}
	clients[c.ID()] = c
}

// unsubscribe drops the client from the topic and the topic itself once empty,
// otherwise an idle topic would leak its bucket array forever.
func (s *EventShard) unsubscribe(clientID, topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients, ok := s.eventSubsMap[topic]
	if !ok {
		return
	}

	delete(clients, clientID)
	if len(clients) == 0 {
		delete(s.eventSubsMap, topic)
	}
}

// subscribers returns a copy so callers can fan out without holding the lock.
func (s *EventShard) subscribers(topic string) []outbound.Connection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := s.eventSubsMap[topic]
	out := make([]outbound.Connection, 0, len(clients))
	for _, c := range clients {
		out = append(out, c)
	}
	return out
}
