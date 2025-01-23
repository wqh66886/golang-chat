package client

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Socket  *websocket.Conn
	Receive chan []byte
	Room    *Room
}

func (c *Client) Read() {
	defer func() {
		err := c.Socket.Close()
		if err != nil {
			log.Println("web socket close error: ", err)
		}
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Println("read message error: ", err)
			return
		}
		c.Room.Forward <- message
	}
}

func (c *Client) Write() {
	defer func() {
		err := c.Socket.Close()
		if err != nil {
			log.Println("web socket close error: ", err)
		}
	}()

	for bytes := range c.Receive {
		err := c.Socket.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Println("write message error:", err)
			return
		}
	}
}
