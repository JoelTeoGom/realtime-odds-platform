package main

import (
	"log"
	"net/http"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/adapters/inbound/websocket"
	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/service/hub"
)

func main() {
	//hub creation with 5 shards + settings
	hub := hub.NewHub(5)

	handlerCfg := websocket.NewHandler(hub, []string{})

	//Endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handlerCfg.ServeHTTP)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("listen:", err)
	}
}
