package dto

import "time"

type Authenticate struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type AuthenticateResponse struct {
	Token     string    `json:"access_token"`
	TokenType string    `json:"token_type"`
	ExpiresIn time.Time `json:"expires_in"`
}

type ChangePassword struct {
	Username        string `json:"-" swaggerignore:"true"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type TokenRequest struct {
	Key  string `json:"token"`
	Type string `json:"type"`
}
