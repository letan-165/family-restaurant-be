package handler

import (
	"log"
	"myapp/module/notification/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func OrdersWSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	models.Notifier.Register(conn)
	defer models.Notifier.Unregister(conn)

	log.Println("Client connected to /ws/orders")

	for {
		var msg interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Client disconnected:", err)
			break
		}

		models.Notifier.Broadcast(msg)
	}
}
