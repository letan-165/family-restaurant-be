package dto

import "myapp/module/order/models"

type ItemOrderSaveRequest struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderSaveRequest struct {
	Customer models.CustomerOrder   `json:"customer"`
	Items    []ItemOrderSaveRequest `json:"items"`
}
