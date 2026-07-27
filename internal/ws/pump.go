package ws

import (
	"context"
	"log"

	"main/internal/domain"
)

func readPump(c *domain.Client) {
	for {
		mt, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
	}
}

func writePump(ctx context.Context, c *domain.Client) {
	for {
		select {
		case <-ctx.Done():
		case msg, ok := <-c.MsgQueue:
			if !ok {

			}
		}
	}

}
