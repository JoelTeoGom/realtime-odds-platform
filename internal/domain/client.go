package domain

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Id       string
	Conn     *websocket.Conn
	MsgQueue chan string
}

func (c *Client) Send(message string) {
	select {
	case c.MsgQueue <- message:
	default: //id rather lose some client message (if its slow conn)
		//instead of blocking the other users
	}
}
