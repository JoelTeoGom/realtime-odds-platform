package hub

import "sync"

type shard struct {
	mu sync.RWMutex
}
