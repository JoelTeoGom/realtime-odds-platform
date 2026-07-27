package main

import (
	"log"
	"net/http"

	"main/internal/hub"
	"main/internal/ws"
)

func main() {
	h := hub.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.Handler(h))

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("listen:", err)
	}
}
