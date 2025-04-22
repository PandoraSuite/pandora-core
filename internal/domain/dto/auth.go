package dto

import "time"

type Authenticate struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type TokenResponse struct {
	Token     string    `json:"access_token"`
	TokenType string    `json:"token_type"`
	ExpiresIn time.Time `json:"expires_in" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

type AuthenticateResponse struct {
	*TokenResponse     `json:",inline"`
	ForcePasswordReset bool `json:"force_password_reset"`
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
