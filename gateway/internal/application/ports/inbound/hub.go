package inbound

import "github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/ports/outbound"

type Hub interface {
	Subscribe(c outbound.Connection, topics ...string)
	Unsubscribe(c outbound.Connection, topics ...string)
}
