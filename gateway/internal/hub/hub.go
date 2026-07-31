package hub

import (
	"context"
	"hash/fnv"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/ports/inbound"
	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/ports/outbound"
)

type Hub struct {
	EventShards []*EventShard
	ClientShard []*ClientShard
	nShards     int
}

var _ inbound.Hub = (*Hub)(nil)

func NewHub(nShards int) *Hub {
	if nShards < 1 {
		nShards = 1
	}

	eventShards := make([]*EventShard, nShards)
	clientShard := make([]*ClientShard, nShards)
	for i := range shards {
		eventShards[i] = newEventShard()
		clientShard[i] = newClientShard()

	}

	h := &Hub{
		nShards:     nShards,
		EventShards: eventShards,
		ClientShard: clientShards,
	}

	// debug only: dumps the shard maps every 10s.
	go h.MonitorHubData(context.Background())

	return h
}

// shardFor always maps a topic to the same shard.
func (h *Hub) shardFor(topic string) *Shard {
	sum := fnv.New32a()
	_, _ = sum.Write([]byte(topic))
	return h.shards[sum.Sum32()%uint32(h.nShards)]
}

func (h *Hub) Subscribe(c outbound.Client, topics ...string) {
	for _, topic := range topics {
		h.shardFor(topic).subscribe(c, topic)
	}
}

func (h *Hub) Unsubscribe(c outbound.Client, topics ...string) {
	id := c.ID()
	for _, topic := range topics {
		h.shardFor(topic).unsubscribe(id, topic)
	}
}
