package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"myapp/module/order/models"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2/google"
)

type fcmMessage struct {
	Message struct {
		Token string            `json:"token,omitempty"`
		Topic string            `json:"topic,omitempty"`
		Data  map[string]string `json:"data"`
	} `json:"message"`
}

func SendFCMBooking(order models.Order) error {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	credPath := os.Getenv("FIREBASE_CREDENTIALS")

	if projectID == "" || credPath == "" {
		return fmt.Errorf("thiếu FIREBASE_PROJECT_ID hoặc FIREBASE_CREDENTIALS trong .env")
	}

	data, err := os.ReadFile(credPath)
	if err != nil {
		return fmt.Errorf("không đọc được file credentials: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return fmt.Errorf("không tạo được JWT config: %v", err)
	}

	token, err := conf.TokenSource(context.Background()).Token()
	if err != nil {
		return fmt.Errorf("không tạo được access token: %v", err)
	}

	var msg fcmMessage
	msg.Message.Topic = "all"
	msg.Message.Data = map[string]string{
		"title":        "Đơn hàng mới!",
		"order_id":     order.ID.Hex(),
		"time_booking": order.TimeBooking.Format("2006-01-02T15:04:05.000Z"),
		"total":        strconv.Itoa(order.Total),
	}

	bodyJSON, _ := json.Marshal(msg)
	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return fmt.Errorf("lỗi tạo request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("lỗi gửi FCM: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM lỗi: %s", resp.Status)
	}

	fmt.Println("✅ Gửi thông báo booking FCM thành công!")
	return nil
}
