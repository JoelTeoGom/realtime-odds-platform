package inbound

import "github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/ports/outbound"

type Hub interface {
	Register(c outbound.Connection, topic ...string)
	UnRegister(c outbound.Connection, topic ...string)
	Subscribe(c outbound.Connection, topics ...string)
	Unsubscribe(c outbound.Connection, topics ...string)
}
