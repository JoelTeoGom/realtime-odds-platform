package hub

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

const monitorInterval = 10 * time.Second

func (h *Hub) MonitorHubData(ctx context.Context) {
	ticker := time.NewTicker(monitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Print(h.dump())
		}
	}
}

func (h *Hub) dump() string {
	var b strings.Builder

	totalTopics, totalSubs := 0, 0
	for i, s := range h.shards {
		topics := s.snapshot()

		names := make([]string, 0, len(topics))
		for topic := range topics {
			names = append(names, topic)
		}
		sort.Strings(names)

		fmt.Fprintf(&b, "  shard[%d] topics=%d\n", i, len(names))
		for _, topic := range names {
			clients := topics[topic]
			sort.Strings(clients)
			fmt.Fprintf(&b, "    %q -> %d client(s) %v\n", topic, len(clients), clients)
			totalSubs += len(clients)
		}
		totalTopics += len(names)
	}

	header := fmt.Sprintf("=== HUB DUMP shards=%d topics=%d subscriptions=%d ===\n",
		h.nShards, totalTopics, totalSubs)

	return header + b.String()
}

func (s *Shard) snapshot() map[string][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string][]string, len(s.subs))
	for topic, clients := range s.subs {
		ids := make([]string, 0, len(clients))
		for id := range clients {
			ids = append(ids, id)
		}
		out[topic] = ids
	}
	return out
}
