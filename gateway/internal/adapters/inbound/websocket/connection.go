package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/JoelTeoGom/go-sharded-ws-hub/gateway/internal/application/ports/inbound"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
	maxMessage = 4 * 1024
	sendBuffer = 256
)

type Connection struct {
	id   string
	conn *websocket.Conn
	send chan []byte

	once sync.Once
	done chan struct{}
}

func newWsConnection(id string, conn *websocket.Conn) *Connection {
	return &Connection{
		id:   id,
		conn: conn,
		send: make(chan []byte, sendBuffer),
		done: make(chan struct{}),
	}
}

type Command struct {
	Action  string `json:"action"`
	EventID string `json:"eventId"`
}

// Close is idempotent and unblocks both pumps.
func (c *Connection) Close() error {
	c.once.Do(func() {
		close(c.done)
		c.conn.Close()
	})
	return nil
}

func (c *Connection) ID() string {
	return c.id
}

func (c *Connection) Send(payload []byte) {
	select {
	case c.send <- payload:
	default: //We can discard some messages NOT BLOCK
	}
}

func (c *Connection) readPump(hub inbound.Hub) {
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
			hub.Subscribe(c, cmd.EventID)
		case "unsubscribe":
			hub.Unsubscribe(c, cmd.EventID)
		}
	}
}

func (c *Connection) writePump(hub inbound.Hub) {
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
