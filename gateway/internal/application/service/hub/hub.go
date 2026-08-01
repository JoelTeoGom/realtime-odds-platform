package hub

import (
	"context"
	"hash/fnv"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/ports/inbound"
	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/ports/outbound"
)

type Hub struct {
	EventShards []*EventShard
	nShards     int
}

var _ inbound.Hub = (*Hub)(nil)

func NewHub(nShards int) *Hub {
	if nShards < 1 {
		nShards = 1
	}

	eventShards := make([]*EventShard, nShards)
	for i := 0; i < nShards; i++ {
		eventShards[i] = newEventShard()
	}

	h := &Hub{
		nShards:     nShards,
		EventShards: eventShards,
	}

	// debug only: dumps the shard maps every 10s.
	go h.MonitorHubData(context.Background())

	return h
}

func (h *Hub) shardFor(topic string) *EventShard {
	sum := fnv.New32a()
	_, _ = sum.Write([]byte(topic))
	return h.EventShards[sum.Sum32()%uint32(h.nShards)]
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
