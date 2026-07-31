package inbound

import "github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/ports/outbound"

type Hub interface {
	Subscribe(c outbound.Client, topics ...string)
	Unsubscribe(c outbound.Client, topics ...string)
}
