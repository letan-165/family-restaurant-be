package ws

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func OrdersWSHandler(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(conn *websocket.Conn) {
		Notifier.Register(conn)
		defer Notifier.Unregister(conn)

		for {
			var msg interface{}
			if err := websocket.JSON.Receive(conn, &msg); err != nil {
				log.Println("Client disconnected:", err)
				break
			}
		}
	}).ServeHTTP(w, r)
}
