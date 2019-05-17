package ws

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients *sync.Map
}

func NewHub() *Hub {
	return &Hub{
		clients: &sync.Map{},
	}
}

func (hub *Hub) listenCloseConnection(id string, conn *websocket.Conn) {
loop:
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			hub.Remove(id)
			break loop
		}
	}
}

func (hub *Hub) Add(id string, w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, _ := upgrader.Upgrade(w, r, nil)
	hub.clients.Store(id, conn)
	go hub.listenCloseConnection(id, conn)
	return conn
}

func (hub *Hub) Get(id string) (*websocket.Conn, bool) {
	val, ok := hub.clients.Load(id)
	if ok {
		return val.(*websocket.Conn), ok
	}
	return nil, ok
}

func (hub *Hub) Write(id, message string) error {
	conn, ok := hub.Get(id)
	if !ok {
		return fmt.Errorf("Connection with id %s is not exists", id)
	}
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (hub *Hub) Remove(id string) {
	if conn, ok := hub.Get(id); ok {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		conn.WriteMessage(websocket.CloseMessage, []byte{})
		conn.Close()
		hub.clients.Delete(id)
	}
}
