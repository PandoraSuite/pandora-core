package dto

import "time"

type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateResponse struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenRequest struct {
	Key  string `json:"token"`
	Type string `json:"type"`
}
