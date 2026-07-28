package gateway

import (
	"log"
	"net/http"

	"github.com/JoelTeoGom/go-sharded-ws-hub/internal/adapters/inbound/ws"
)

func main() {
	h := hub.New()

	handlerCfg := ws.NewHandler(h)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handlerCfg.WsHandler())

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("listen:", err)
	}
}
