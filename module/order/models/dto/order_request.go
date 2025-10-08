package dto

import "myapp/module/order/models"

type ItemOrderCreateRequest struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderCreateRequest struct {
	Customer models.CustomerOrder     `json:"customer"`
	Items    []ItemOrderCreateRequest `json:"items"`
}
