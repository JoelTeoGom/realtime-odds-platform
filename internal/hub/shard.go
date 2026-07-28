package hub

import (
	"sync"

	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/ports/outbound"
)

type Shard struct {
	mu        sync.Mutex
	clientMap map[string][]*outbound.Client
}
