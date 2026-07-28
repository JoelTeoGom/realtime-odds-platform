package main

import (
	"log"
	"net/http"

	"main/internal/hub"
	"main/internal/ws"
)

func main() {
	h := hub.New()
	h := 
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.Handler(h))
	mux.HandleFunc("/event/update", )

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("listen:", err)
	}

}
