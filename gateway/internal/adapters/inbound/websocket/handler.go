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

	raw, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := newWsConnection(id, raw)

	if err := h.hub.Register(conn); err != nil {
		_ = conn.Close()
		_ = raw.Close()
		h.log.Warn("register rejected", "client", id, "err", err)
		return
	}

	// El writePump arranca primero: garantiza que cualquier Send() disparado
	// por Register (p.ej. mensajes de bienvenida) tenga ya consumidor.
	go conn.writePump()
	go conn.readPump(h.hub, h.log)
}
