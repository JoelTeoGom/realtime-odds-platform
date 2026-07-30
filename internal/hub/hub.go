package hub

import (
	"context"
	"hash/fnv"

	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/ports/inbound"
	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/ports/outbound"
)

type Hub struct {
	shards  []*Shard
	nShards int
}

var _ inbound.Hub = (*Hub)(nil)

func NewHub(nShards int) *Hub {
	if nShards < 1 {
		nShards = 1
	}

	shards := make([]*Shard, nShards)
	for i := range shards {
		shards[i] = newShard()
	}

	h := &Hub{
		nShards: nShards,
		shards:  shards,
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
