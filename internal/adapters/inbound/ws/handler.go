package ws

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/ports/inbound"
)

type Handler struct {
	hub      inbound.Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub inbound.Hub) *Handler {
	return &Handler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// TODO: restrict to your own origins before going public.
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) WsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		c := newClient(conn, h.hub)

		// no defer here: this func returns right away and the connection
		// keeps living in the pumps. readPump is the one that unregisters.
		go c.writePump()
		go c.readPump()
	}
}
