package domain

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id       string
	Conn     *websocket.Conn
	MsgQueue chan Event
}

func NewClient(conn *websocket.Conn, msqQueue chan Event) *Client {
	return &Client{
		Id:       uuid.NewString(),
		Conn:     conn,
		MsgQueue: msqQueue,
	}
}

func (c *Client) Send(message Event) {
	select {
	case c.MsgQueue <- message:
	default:
		//id rather lose some client message (if its slow conn)
		//instead of blocking the other users
	}
}
