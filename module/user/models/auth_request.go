package models

type GenerateTokenAdminRequest struct {
	Email  string `json:"email"`
	Secret string `json:"secret"`
}
