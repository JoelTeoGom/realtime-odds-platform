package websocket

import (
	"net/http"
	"slices"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/ports/inbound"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub      inbound.Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub inbound.Hub, allowedOrigins []string) *Handler {
	return &Handler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return slices.Contains(allowedOrigins, r.Header.Get("Origin"))
			},
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("client_id") //TODO: prod will have JWT (for testing purposes we take it from query params)
	if id == "" {
		http.Error(w, "missing client_id", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	wsConn := newWsConnection(id, conn)
	defer func() {
		wsConn.Close()
	}()

	if err := h.hub.Register(conn); err != nil {
		return
	}

	go wsConn.readPump(h.hub)
	wsConn.writePump(h.hub)
}
