package services

import "myapp/module/notification/ws"

func SendNewOrderNotification(orderId int, status string) {
	msg := map[string]interface{}{
		"orderId": orderId,
		"status":  status,
	}
	ws.Notifier.Broadcast(msg)
}
