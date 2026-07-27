package hub

import (
	"sync"

	"main/internal/domain"
)

type Hub struct {
	Mu        sync.RWMutex
	ClientMap map[string][]*domain.Client
}

func New() *Hub {
	return &Hub{
		ClientMap: make(map[string][]*domain.Client),
	}
}

func (h *Hub) Register(key string, value *domain.Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	h.ClientMap[key] = append(h.ClientMap[key], value)
}

func (h *Hub) Unregister(key string, value *domain.Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	list := h.ClientMap[key]
	for i, client := range list {
		if client == value {
			list[i] = list[len(list)-1]
			list[len(list)-1] = nil
			list = list[:len(list)-1]
			break
		}
	}

	if len(list) == 0 {
		delete(h.ClientMap, key)
		return
	}

	h.ClientMap[key] = list

}

func (h *Hub) Broadcast() {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

}
