package client

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Room struct {
	Forward chan []byte
	Clients map[*Client]bool
	Join    chan *Client
	Leave   chan *Client
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan []byte),
		Clients: make(map[*Client]bool),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {
		case cli := <-r.Join:
			r.Clients[cli] = true
		case cli := <-r.Leave:
			delete(r.Clients, cli)
			close(cli.Receive) //关闭接收消息的channel
		case bytes := <-r.Forward:
			for cli := range r.Clients {
				cli.Receive <- bytes
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("web socket upgrade error: ", err)
		return
	}
	client := &Client{
		Receive: make(chan []byte, messageBufferSize),
		Room:    r,
		Socket:  socket,
	}
	r.Join <- client
	defer func() {
		r.Leave <- client
	}()
	go client.Write()
	client.Read()
}
