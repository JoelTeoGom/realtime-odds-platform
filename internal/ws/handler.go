package ws

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"main/internal/domain"
	"main/internal/hub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := &domain.Client{
			Id: uuid.NewString(),
		}
	}
}
