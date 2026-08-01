package hub

import (
	"context"
	"hash/fnv"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/ports/inbound"
	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/ports/outbound"
)

type Hub struct {
	EventShards []*EventShard
	UserShard   []*UserShard
	nShards     int
}

var _ inbound.Hub = (*Hub)(nil)

func NewHub(nShards int) *Hub {
	if nShards < 1 {
		nShards = 1
	}

	eventShards := make([]*EventShard, nShards)
	userShards := make([]*UserShard, nShards)
	for i := 0; i < nShards; i++ {
		eventShards[i] = newEventShard()
		userShards[i] = newUserShard()
	}

	h := &Hub{
		nShards:     nShards,
		EventShards: eventShards,
		UserShard:   userShards,
	}

	// debug only: dumps the shard maps every 10s.
	go h.MonitorHubData(context.Background())

	return h
}

// shardFor always maps a topic to the same shard.
func (h *Hub) shardFor(topic string) *UserShard {
	sum := fnv.New32a()
	_, _ = sum.Write([]byte(topic))
	return h.UserShard[sum.Sum32()%uint32(h.nShards)]
}

func (h *Hub) Subscribe(c outbound.Connection, topics ...string) {
	for _, topic := range topics {
		h.shardFor(topic).subscribe(c, topic)
	}
}

func (h *Hub) Unsubscribe(c outbound.Connection, topics ...string) {
	id := c.ID()
	for _, topic := range topics {
		h.shardFor(topic).unsubscribe(id, topic)
	}
}
