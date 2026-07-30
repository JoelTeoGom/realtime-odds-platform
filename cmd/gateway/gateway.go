package gateway

import (
	"log"
	"net/http"

	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/adapters/inbound/ws"
	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/hub"
)

func NewGateway() {
	//hub creation with 5 shards + settings
	hub := hub.NewHub(5)
	handlerCfg := ws.NewHandler(hub)

	//Endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handlerCfg.WsHandler())

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("listen:", err)
	}
}
