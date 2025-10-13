package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Notification struct {
	Clients map[*websocket.Conn]bool
	Mu      sync.Mutex
}

var Notifier = &Notification{
	Clients: make(map[*websocket.Conn]bool),
}

func (n *Notification) Register(conn *websocket.Conn) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	n.Clients[conn] = true
}

func (n *Notification) Broadcast(msg interface{}) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	for conn := range n.Clients {
		if err := conn.WriteJSON(msg); err != nil {
			conn.Close()
			delete(n.Clients, conn)
		}
	}
}

func (n *Notification) Unregister(conn *websocket.Conn) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	conn.Close()
	delete(n.Clients, conn)
}
