package dto

import "time"

type Authenticate struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type AuthenticateResponse struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ChangePassword struct {
	Username        string `json:"-"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type TokenRequest struct {
	Key  string `json:"token"`
	Type string `json:"type"`
}
