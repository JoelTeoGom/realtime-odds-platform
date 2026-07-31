package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/ports/inbound"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
	maxMessage = 4 * 1024
	sendBuffer = 256
)

type Client struct {
	id   string
	conn *websocket.Conn
	hub  inbound.Hub

	send chan []byte
	once sync.Once
	done chan struct{} // closed by Close, unblocks writePump
}

type Command struct {
	Action  string `json:"action"`
	EventID string `json:"eventId"`
}

func newClient(conn *websocket.Conn, hub inbound.Hub) *Client {
	return &Client{
		id:   uuid.NewString(),
		conn: conn,
		hub:  hub,
		send: make(chan []byte, sendBuffer),
		done: make(chan struct{}),
	}
}

// Close is idempotent and unblocks both pumps.
func (c *Client) Close() error {
	c.once.Do(func() {
		close(c.done)
		_ = c.conn.Close()
	})
	return nil
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Send(payload []byte) {
	select {
	case c.send <- payload:
	default: //We can discard some messages NOT BLOCK
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Close()
	}()

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return //if error close
		}

		var cmd Command
		if err := json.Unmarshal(raw, &cmd); err != nil {
			continue //trash we ignore
		}

		switch cmd.Action {
		case "subscribe":
			c.hub.Subscribe(c, cmd.EventID)
		case "unsubscribe":
			c.hub.Unsubscribe(c, cmd.EventID)
		}
	}
}

func (c *Client) writePump() {
	ping := time.NewTicker(pingPeriod)
	defer func() {
		ping.Stop()
		c.Close()
	}()

	for {
		select {
		case payload := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}

		case <-ping.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:
			return //unblock the write
		}
	}
}
